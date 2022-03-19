package payment_gateway

import (
	"context"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/overridesh/checkout-challenger/internal/model"
	"github.com/overridesh/checkout-challenger/internal/repository"
	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
	"github.com/overridesh/checkout-challenger/pkg/storage/cache"
	pbPaymentGateway "github.com/overridesh/checkout-challenger/proto"
	"github.com/overridesh/checkout-challenger/tools"
)

// NewPaymentGatewayGRPC implements the protobuf interface
type PaymentGatewayGRPC struct {
	repository    repository.PaymentGatewayRepository
	bankSimulator acquirer.Acquirer
	cache         cache.Cache
}

// New initializes a new NewPaymentGatewayGRPC struct.
func NewGRPC(repository repository.PaymentGatewayRepository, bankSimulatorService acquirer.Acquirer, cache cache.Cache) *PaymentGatewayGRPC {
	return &PaymentGatewayGRPC{
		repository:    repository,
		bankSimulator: bankSimulatorService,
		cache:         cache,
	}
}

func (pg *PaymentGatewayGRPC) CreatePayment(ctx context.Context, in *pbPaymentGateway.CreatePaymentRequest) (*pbPaymentGateway.CreatePaymentResponse, error) {
	merchantID, err := GetMerchantID(ctx)
	if err != nil {
		return nil, err
	}

	var (
		idempotencyItem model.Idempotency
		step            model.RecoveryPoint = model.FirstPoint
	)

	idempotencyKey, ok := GetIdempotencyKey(ctx)
	if ok && len(idempotencyKey) > 0 {
		item, itemOk := pg.cache.Get(idempotencyKey)
		if itemOk {
			idempotencyItem, itemOk = item.(model.Idempotency)
			if itemOk {
				if idempotencyItem.RecoveryPoint == model.UpdatedPoint {
					return idempotencyItem.Response, nil
				}
				step = idempotencyItem.RecoveryPoint
			}
		}
	}

	var (
		rawCard *pbPaymentGateway.CreditCard = in.GetCreditCard()
	)

	if err := tools.CardValidate(rawCard.FirstName, rawCard.LastName, rawCard.Number, rawCard.Cvv, rawCard.ExpiryMonth, rawCard.ExpiryYear); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if step == model.FirstPoint {
		return pg.FromFirstPoint(ctx, model.Transaction{
			MerchantID:        merchantID,
			Amount:            in.GetAmount(),
			Currency:          in.GetCurrency(),
			Reference:         in.GetReference(),
			SourceNumber:      in.CreditCard.Number,
			SourceFirstName:   in.CreditCard.GetFirstName(),
			SourceLastName:    in.CreditCard.GetLastName(),
			SourceExpiryMonth: in.CreditCard.GetExpiryMonth(),
			SourceExpiryYear:  in.CreditCard.GetExpiryYear(),
			IdempotencyKey:    idempotencyKey,
		})
	}

	transaction, err := pg.repository.GetByMerchantIDAndIdempotencyKey(merchantID, idempotencyKey)
	if err != nil {
		if errors.Is(repository.ErrTransactionNotFound, err) {
			return pg.FromFirstPoint(ctx, model.Transaction{
				MerchantID:        merchantID,
				Amount:            in.GetAmount(),
				Currency:          in.GetCurrency(),
				Reference:         in.GetReference(),
				SourceNumber:      in.CreditCard.Number,
				SourceFirstName:   in.CreditCard.GetFirstName(),
				SourceLastName:    in.CreditCard.GetLastName(),
				SourceExpiryMonth: in.CreditCard.GetExpiryMonth(),
				SourceExpiryYear:  in.CreditCard.GetExpiryYear(),
				IdempotencyKey:    idempotencyKey,
			})
		}
		return nil, err
	}

	switch step {
	case model.CreatedPoint:
		return pg.FromCreatedTrasactionPoint(ctx, transaction)
	case model.PurchasedPoint:
		return pg.FromPuchasedTransactionPoint(ctx, transaction, &idempotencyItem)
	case model.UpdatedPoint:
		return pg.SuccessTransactionStep(transaction), nil
	}

	return nil, tools.SetStatusCode(ctx, http.StatusUnprocessableEntity)
}

func (pg *PaymentGatewayGRPC) GetPayment(ctx context.Context, in *pbPaymentGateway.GetPaymentDetailsRequest) (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
	merchantID, err := GetMerchantID(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromString(in.Id)
	if err != nil {
		return nil, IdMustBeUUID.Err()
	}

	if id == uuid.Nil {
		return nil, IdMusBeValidUUID.Err()
	}

	payment, err := pg.repository.GetByID(id, merchantID)
	if err != nil {
		if errors.Is(err, repository.ErrTransactionNotFound) {
			return nil, transactionNotFound.Err()
		}

		zap.S().Errorw("error get transaction by id", zap.Any("payload", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}))

		return nil, internalServerErr.Err()
	}

	lastFour, _ := tools.LastFour(payment.SourceNumber)

	return &pbPaymentGateway.GetPaymentDetailsResponse{
		Id:       payment.Id.String(),
		Amount:   payment.Amount,
		Currency: payment.Currency,
		Source: &pbPaymentGateway.Source{
			FirstName:   payment.SourceFirstName,
			LastName:    payment.SourceLastName,
			Last4:       lastFour,
			Bin:         payment.SourceBin.String,
			CardType:    payment.SourceCardType.String,
			ExpiryMonth: payment.SourceExpiryMonth,
			ExpiryYear:  payment.SourceExpiryYear,
		},
		Status:      payment.Status.String(),
		Approved:    payment.Approved,
		RequestedOn: tools.FormatDate(payment.CreatedAt),
	}, nil
}
