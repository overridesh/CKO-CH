package repository

import (
	"database/sql"
	"errors"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/overridesh/checkout-challenger/internal/model"
	storage "github.com/overridesh/checkout-challenger/pkg/storage/sql"
)

var (
	ErrTransactionNotFound = errors.New("payment not found")
)

const (
	getByID string = `
	 	SELECT 
		 id, 
		 merchant_id, 
		 approved, 
		 status, 
		 amount, 
		 currency, 
		 source_first_name, 
		 source_last_name, 
		 source_number, 
		 source_bin, 
		 source_card_type, 
		 source_expiry_month, 
		 source_expiry_year, 
		 response_code, 
		 response_summary, 
		 reference,
		 created_at,
		 idempotency_key
	 	FROM transactions 
	 	WHERE id = $1
		AND merchant_id = $2
	`
	getByMerchantIdAndIdempotencyKey string = `
	 	SELECT 
		 id, 
		 merchant_id, 
		 approved, 
		 status, 
		 amount, 
		 currency, 
		 source_first_name, 
		 source_last_name, 
		 source_number, 
		 source_bin, 
		 source_card_type, 
		 source_expiry_month, 
		 source_expiry_year, 
		 response_code, 
		 response_summary, 
		 reference,
		 created_at,
		 idempotency_key
	 	FROM transactions 
	 	WHERE merchant_id = $1
		AND idempotency_key = $2
	`

	updateByID string = `
		UPDATE transactions
		 SET 
		 	status = $1,
			response_code = $2,
			response_summary = $3,
			approved = $4,
			source_bin = $5,
			source_card_type = $6,
			source_number = $7,
			source_first_name = $8, 
			source_last_name = $9
		WHERE id = $10
	`
	insertTransaction string = `
		INSERT INTO transactions (
			merchant_id,
			amount,
			currency, 
			source_first_name,
			source_last_name,
			source_expiry_month,
			source_expiry_year,
			source_number,
			reference,
			idempotency_key
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10
		) RETURNING 
			id,
			merchant_id, 
			approved, 
			status, 
			amount, 
			currency, 
			source_first_name, 
			source_last_name, 
			source_number, 
			source_bin, 
			source_card_type, 
			source_expiry_month, 
			source_expiry_year, 
			response_code, 
			response_summary, 
			reference,
			created_at,
			idempotency_key;
	`
)

type PaymentGatewayRepository interface {
	Create(payment model.Transaction) (*model.Transaction, error)
	Update(payment *model.Transaction, fn func(*model.Transaction) error) error
	GetByID(id, merchantID uuid.UUID) (*model.Transaction, error)
	GetByMerchantIDAndIdempotencyKey(merchantID uuid.UUID, idempotencyKey string) (*model.Transaction, error)
}

type paymentGateway struct {
	db storage.DB
}

func NewPaymentGatewayRepository(db storage.DB) PaymentGatewayRepository {
	return &paymentGateway{
		db: db,
	}
}

func (pg *paymentGateway) Create(payment model.Transaction) (*model.Transaction, error) {
	var newTransaction model.Transaction

	if err := pg.db.QueryRow(
		insertTransaction,
		payment.MerchantID,
		payment.Amount,
		payment.Currency,
		payment.SourceFirstName,
		payment.SourceLastName,
		payment.SourceExpiryMonth,
		payment.SourceExpiryYear,
		payment.SourceNumber,
		payment.Reference,
		payment.IdempotencyKey,
	).Scan(
		&newTransaction.Id,
		&newTransaction.MerchantID,
		&newTransaction.Approved,
		&newTransaction.Status,
		&newTransaction.Amount,
		&newTransaction.Currency,
		&newTransaction.SourceFirstName,
		&newTransaction.SourceLastName,
		&newTransaction.SourceNumber,
		&newTransaction.SourceBin,
		&newTransaction.SourceCardType,
		&newTransaction.SourceExpiryMonth,
		&newTransaction.SourceExpiryYear,
		&newTransaction.ResponseCode,
		&newTransaction.ResponseSummary,
		&newTransaction.Reference,
		&newTransaction.CreatedAt,
		&newTransaction.IdempotencyKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return &newTransaction, nil
}

func (pg *paymentGateway) Update(payment *model.Transaction, tryPayment func(*model.Transaction) error) error {
	var (
		err error
		tx  *sql.Tx
	)

	tx, err = pg.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				zap.S().Errorf("cannot do a rollback, error: %v", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				zap.S().Errorf("cannot do a commit, error: %v", err)
			}
		}
	}()

	errPayment := tryPayment(payment)

	row, err := tx.Query(
		updateByID,
		payment.Status,
		payment.ResponseCode,
		payment.ResponseSummary,
		payment.Approved,
		payment.SourceBin,
		payment.SourceCardType,
		payment.SourceNumber,
		payment.SourceFirstName,
		payment.SourceLastName,
		payment.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrTransactionNotFound
		}
		return err
	}

	row.Close()

	return errPayment
}

func (pg *paymentGateway) GetByID(id, merchantID uuid.UUID) (*model.Transaction, error) {
	var payment model.Transaction
	if err := pg.db.QueryRow(
		getByID,
		id,
		merchantID,
	).Scan(
		&payment.Id,
		&payment.MerchantID,
		&payment.Approved,
		&payment.Status,
		&payment.Amount,
		&payment.Currency,
		&payment.SourceFirstName,
		&payment.SourceLastName,
		&payment.SourceNumber,
		&payment.SourceBin,
		&payment.SourceCardType,
		&payment.SourceExpiryMonth,
		&payment.SourceExpiryYear,
		&payment.ResponseCode,
		&payment.ResponseSummary,
		&payment.Reference,
		&payment.CreatedAt,
		&payment.IdempotencyKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (pg *paymentGateway) GetByMerchantIDAndIdempotencyKey(merchantID uuid.UUID, idempotencyKey string) (*model.Transaction, error) {
	var payment model.Transaction

	if err := pg.db.QueryRow(
		getByMerchantIdAndIdempotencyKey,
		merchantID,
		idempotencyKey,
	).Scan(
		&payment.Id,
		&payment.MerchantID,
		&payment.Approved,
		&payment.Status,
		&payment.Amount,
		&payment.Currency,
		&payment.SourceFirstName,
		&payment.SourceLastName,
		&payment.SourceNumber,
		&payment.SourceBin,
		&payment.SourceCardType,
		&payment.SourceExpiryMonth,
		&payment.SourceExpiryYear,
		&payment.ResponseCode,
		&payment.ResponseSummary,
		&payment.Reference,
		&payment.CreatedAt,
		&payment.IdempotencyKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return &payment, nil
}
