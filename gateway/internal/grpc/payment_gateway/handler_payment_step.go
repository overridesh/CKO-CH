package payment_gateway

import (
	"context"

	"go.uber.org/zap"

	"github.com/overridesh/checkout-challenger/internal/model"
	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
	"github.com/overridesh/checkout-challenger/pkg/storage/cache"
	pbPaymentGateway "github.com/overridesh/checkout-challenger/proto"
	"github.com/overridesh/checkout-challenger/tools"
)

func (pg *PaymentGatewayGRPC) SetRecoveryPoint(key string, item *model.Idempotency) {
	if pg.cache != nil && len(key) > 0 {
		if err := item.SetMD5(); err != nil {
			zap.S().Errorf("error with idempotency marshal, error: %v", err)
			return
		}
		pg.cache.Set(key, *item, cache.DefaultExpiration)
	}
}

// FromFirstPoint recovery from first step
func (pg *PaymentGatewayGRPC) FromFirstPoint(ctx context.Context, tx model.Transaction) (*pbPaymentGateway.CreatePaymentResponse, error) {
	payment, err := pg.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := pg.UpdateTransactionStep(ctx, payment); err != nil {
		return nil, err
	}

	return pg.SuccessTransactionStep(payment), nil
}

// FromCreatedTrasactionPoint Recovery when the transaction was created
func (pg *PaymentGatewayGRPC) FromCreatedTrasactionPoint(ctx context.Context, payment *model.Transaction) (*pbPaymentGateway.CreatePaymentResponse, error) {
	if err := pg.UpdateTransactionStep(ctx, payment); err != nil {
		return nil, err
	}

	return pg.SuccessTransactionStep(payment), nil
}

// FromPuchasedTransactionPoint Recovery when the transaction was paided
func (pg *PaymentGatewayGRPC) FromPuchasedTransactionPoint(ctx context.Context, payment *model.Transaction, idemItem *model.Idempotency) (*pbPaymentGateway.CreatePaymentResponse, error) {
	if err := pg.repository.Update(payment, func(tx *model.Transaction) error {
		tx.SourceBin.String = idemItem.BankResponse.CardBin
		tx.SourceBin.Valid = true

		tx.SourceCardType.String = idemItem.BankResponse.CardType
		tx.SourceCardType.Valid = true

		tx.Status = model.NewStatus(idemItem.BankResponse.Status)
		tx.Approved = tx.Status.IsApproved()

		tx.ResponseCode.String = idemItem.BankResponse.Code
		tx.ResponseCode.Valid = true

		tx.ResponseSummary.String = idemItem.BankResponse.Summary
		tx.ResponseSummary.Valid = true

		tools.SetStatusCode(ctx, idemItem.BankResponse.StatusCode)
		return nil
	}); err != nil {
		return nil, err
	}

	return pg.SuccessTransactionStep(payment), nil
}

func (pg *PaymentGatewayGRPC) CreateTransaction(ctx context.Context, tx model.Transaction) (*model.Transaction, error) {
	item := model.Idempotency{
		RecoveryPoint: model.FirstPoint,
		Response:      nil,
		MerchantID:    tx.MerchantID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		FirstName:     tx.SourceFirstName,
		LastName:      tx.SourceLastName,
		Number:        tx.SourceNumber,
		ExpiryMonth:   tx.SourceExpiryMonth,
		ExpiryYear:    tx.SourceExpiryYear,
	}

	defer pg.SetRecoveryPoint(tx.IdempotencyKey, &item)

	payment, err := pg.repository.Create(tx)
	if err != nil {
		zap.S().Errorw("error create transaction", zap.Any("payload", map[string]interface{}{
			"error":       err.Error(),
			"merchant_id": tx.MerchantID,
			"amount":      tx.Amount,
			"currency":    tx.Currency,
		}))

		return nil, err
	}

	item.RecoveryPoint = model.CreatedPoint

	return payment, nil
}

func (pg *PaymentGatewayGRPC) UpdateTransactionStep(ctx context.Context, tx *model.Transaction) error {
	item := model.Idempotency{
		RecoveryPoint: model.CreatedPoint,
		MerchantID:    tx.MerchantID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		FirstName:     tx.SourceFirstName,
		LastName:      tx.SourceLastName,
		Number:        tx.SourceNumber,
		ExpiryMonth:   tx.SourceExpiryMonth,
		ExpiryYear:    tx.SourceExpiryYear,
	}

	defer pg.SetRecoveryPoint(tx.IdempotencyKey, &item)

	return pg.repository.Update(tx, func(tx *model.Transaction) error {
		response, err := pg.bankSimulator.Purchase(&acquirer.PaymentRequest{
			Amount:      tx.Amount,
			Currency:    tx.Currency,
			Number:      tx.SourceNumber,
			FirstName:   tx.SourceFirstName,
			LastName:    tx.SourceLastName,
			ExpiryMonth: tx.SourceExpiryMonth,
			ExpiryYear:  tx.SourceExpiryYear,
		})

		if err != nil && response == nil {
			return err
		}

		status := model.NewStatus(response.Status)
		if status.IsApproved() {
			item.RecoveryPoint = model.PurchasedPoint
			item.BankResponse = response
		}

		tx.SourceBin.String = response.CardBin
		tx.SourceBin.Valid = true

		tx.SourceCardType.String = response.CardType
		tx.SourceCardType.Valid = true

		tx.Status = status
		tx.Approved = tx.Status.IsApproved()

		tx.ResponseCode.String = response.Code
		tx.ResponseCode.Valid = true

		tx.ResponseSummary.String = response.Summary
		tx.ResponseSummary.Valid = true

		tx.SourceFirstName = response.FirstName
		tx.SourceLastName = response.LastName

		tools.SetStatusCode(ctx, response.StatusCode)
		return err
	})
}

func (pg *PaymentGatewayGRPC) SuccessTransactionStep(tx *model.Transaction) *pbPaymentGateway.CreatePaymentResponse {
	lastFour, _ := tools.LastFour(tx.SourceNumber)

	response := &pbPaymentGateway.CreatePaymentResponse{
		Id:       tx.Id.String(),
		Amount:   tx.Amount,
		Currency: tx.Currency,
		Source: &pbPaymentGateway.Source{
			FirstName:   tx.SourceFirstName,
			LastName:    tx.SourceLastName,
			Last4:       lastFour,
			Bin:         tx.SourceBin.String,
			CardType:    tx.SourceCardType.String,
			ExpiryMonth: tx.SourceExpiryMonth,
			ExpiryYear:  tx.SourceExpiryYear,
		},
		Status:       tx.Status.String(),
		Approved:     tx.Approved,
		ProcessedOn:  tools.FormatDate(tx.CreatedAt),
		ResponseCode: tx.ResponseCode.String,
		Reference:    tx.Reference,
	}

	pg.SetRecoveryPoint(tx.IdempotencyKey, &model.Idempotency{
		RecoveryPoint: model.UpdatedPoint,
		Response:      response,
		MerchantID:    tx.MerchantID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		FirstName:     tx.SourceFirstName,
		LastName:      tx.SourceLastName,
		Number:        tx.SourceNumber,
		ExpiryMonth:   tx.SourceExpiryMonth,
		ExpiryYear:    tx.SourceExpiryYear,
	})

	return &pbPaymentGateway.CreatePaymentResponse{
		Id:       tx.Id.String(),
		Amount:   tx.Amount,
		Currency: tx.Currency,
		Source: &pbPaymentGateway.Source{
			FirstName:   tx.SourceFirstName,
			LastName:    tx.SourceLastName,
			Last4:       lastFour,
			Bin:         tx.SourceBin.String,
			CardType:    tx.SourceCardType.String,
			ExpiryMonth: tx.SourceExpiryMonth,
			ExpiryYear:  tx.SourceExpiryYear,
		},
		Status:       tx.Status.String(),
		Approved:     tx.Approved,
		ProcessedOn:  tools.FormatDate(tx.CreatedAt),
		ResponseCode: tx.ResponseCode.String,
		Reference:    tx.Reference,
	}
}
