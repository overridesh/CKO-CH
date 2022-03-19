package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"

	"github.com/overridesh/checkout-challenger/internal/model"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*model.Transaction, error)
		expect error
	}{
		{
			name: "Create_ErrorNoRows",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				payment := model.Transaction{}

				mock.ExpectQuery(regexp.QuoteMeta(insertTransaction)).WithArgs(
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
				).WillReturnError(sql.ErrNoRows)

				svc := NewPaymentGatewayRepository(db)

				return svc.Create(model.Transaction{})
			},
			expect: ErrTransactionNotFound,
		},
		{
			name: "Create_ErrorConnDone",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				payment := model.Transaction{}

				mock.ExpectQuery(regexp.QuoteMeta(insertTransaction)).WithArgs(
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
				).WillReturnError(sql.ErrConnDone)

				svc := NewPaymentGatewayRepository(db)

				return svc.Create(model.Transaction{})
			},
			expect: sql.ErrConnDone,
		},
		{
			name: "Create_Success",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				payment := model.Transaction{
					Id:                uuid.NewV4(),
					MerchantID:        uuid.NewV4(),
					Amount:            100,
					Currency:          "USD",
					SourceFirstName:   "John",
					SourceLastName:    "Doe",
					SourceExpiryMonth: "12",
					SourceExpiryYear:  "2222",
					SourceNumber:      "4445555444455554444",
				}

				mock.ExpectQuery(regexp.QuoteMeta(insertTransaction)).WithArgs(
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
				).WillReturnRows(sqlmock.NewRows(
					[]string{
						"id",
						"merchant_id",
						"approved",
						"status",
						"amount",
						"currency",
						"source_first_name",
						"source_last_name",
						"source_number",
						"source_bin",
						"source_card_type",
						"source_expiry_month",
						"source_expiry_year",
						"response_code",
						"response_summary",
						"reference",
						"created_at",
						"idempotency_key",
					},
				).AddRow(
					payment.Id,
					payment.MerchantID,
					payment.Approved,
					payment.Status,
					payment.Amount,
					payment.Currency,
					payment.SourceFirstName,
					payment.SourceLastName,
					payment.SourceNumber,
					payment.SourceBin,
					payment.SourceCardType,
					payment.SourceExpiryMonth,
					payment.SourceExpiryYear,
					payment.ResponseCode,
					payment.ResponseSummary,
					payment.Reference,
					payment.CreatedAt,
					payment.IdempotencyKey,
				))

				svc := NewPaymentGatewayRepository(db)

				return svc.Create(payment)
			},
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*model.Transaction, error)
		expect error
	}{
		{
			name: "GetById_ErrorNoRows",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				var (
					merchantId    uuid.UUID = uuid.NewV4()
					transactionId uuid.UUID = uuid.NewV4()
				)
				mock.ExpectQuery(regexp.QuoteMeta(getByID)).WithArgs(
					merchantId,
					transactionId,
				).WillReturnError(sql.ErrNoRows)

				svc := NewPaymentGatewayRepository(db)

				return svc.GetByID(merchantId, transactionId)
			},
			expect: ErrTransactionNotFound,
		},
		{
			name: "GetById_ErrorConnDone",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				var (
					merchantId    uuid.UUID = uuid.NewV4()
					transactionId uuid.UUID = uuid.NewV4()
				)
				mock.ExpectQuery(regexp.QuoteMeta(getByID)).WithArgs(
					merchantId,
					transactionId,
				).WillReturnError(sql.ErrConnDone)

				svc := NewPaymentGatewayRepository(db)

				return svc.GetByID(merchantId, transactionId)
			},
			expect: sql.ErrConnDone,
		},
		{
			name: "GetById_Success",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				payment := model.Transaction{
					Id:                uuid.NewV4(),
					MerchantID:        uuid.NewV4(),
					Amount:            100,
					Currency:          "USD",
					SourceFirstName:   "John",
					SourceLastName:    "Doe",
					SourceExpiryMonth: "12",
					SourceExpiryYear:  "2222",
					SourceNumber:      "4445555444455554444",
				}

				mock.ExpectQuery(regexp.QuoteMeta(getByID)).WithArgs(
					payment.Id,
					payment.MerchantID,
				).WillReturnRows(sqlmock.NewRows(
					[]string{
						"id",
						"merchant_id",
						"approved",
						"status",
						"amount",
						"currency",
						"source_first_name",
						"source_last_name",
						"source_number",
						"source_bin",
						"source_card_type",
						"source_expiry_month",
						"source_expiry_year",
						"response_code",
						"response_summary",
						"reference",
						"created_at",
						"idempotency_key",
					},
				).AddRow(
					payment.Id,
					payment.MerchantID,
					payment.Approved,
					payment.Status,
					payment.Amount,
					payment.Currency,
					payment.SourceFirstName,
					payment.SourceLastName,
					payment.SourceNumber,
					payment.SourceBin,
					payment.SourceCardType,
					payment.SourceExpiryMonth,
					payment.SourceExpiryYear,
					payment.ResponseCode,
					payment.ResponseSummary,
					payment.Reference,
					payment.CreatedAt,
					payment.IdempotencyKey,
				))

				svc := NewPaymentGatewayRepository(db)

				return svc.GetByID(payment.Id, payment.MerchantID)
			},
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestGetByMerchantIDAndIdempotencyKey(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*model.Transaction, error)
		expect error
	}{
		{
			name: "GetByMerchantIDAndIdempotencyKey_ErrorNoRows",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				var (
					merchantId     uuid.UUID = uuid.NewV4()
					idempotencyKey uuid.UUID = uuid.NewV4()
				)
				mock.ExpectQuery(regexp.QuoteMeta(getByMerchantIdAndIdempotencyKey)).WithArgs(
					merchantId,
					idempotencyKey,
				).WillReturnError(sql.ErrNoRows)

				svc := NewPaymentGatewayRepository(db)

				return svc.GetByMerchantIDAndIdempotencyKey(merchantId, idempotencyKey.String())
			},
			expect: ErrTransactionNotFound,
		},
		{
			name: "GetByMerchantIDAndIdempotencyKey_ErrorConnDone",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				var (
					merchantId     uuid.UUID = uuid.NewV4()
					idempotencyKey uuid.UUID = uuid.NewV4()
				)
				mock.ExpectQuery(regexp.QuoteMeta(getByMerchantIdAndIdempotencyKey)).WithArgs(
					merchantId,
					idempotencyKey,
				).WillReturnError(sql.ErrConnDone)

				svc := NewPaymentGatewayRepository(db)

				return svc.GetByMerchantIDAndIdempotencyKey(merchantId, idempotencyKey.String())
			},
			expect: sql.ErrConnDone,
		},
		{
			name: "GetByMerchantIDAndIdempotencyKey_Success",
			input: func() (*model.Transaction, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				payment := model.Transaction{
					Id:                uuid.NewV4(),
					MerchantID:        uuid.NewV4(),
					Amount:            100,
					Currency:          "USD",
					SourceFirstName:   "John",
					SourceLastName:    "Doe",
					SourceExpiryMonth: "12",
					SourceExpiryYear:  "2222",
					SourceNumber:      "4445555444455554444",
					IdempotencyKey:    uuid.NewV4().String(),
				}

				mock.ExpectQuery(regexp.QuoteMeta(getByMerchantIdAndIdempotencyKey)).WithArgs(
					payment.MerchantID,
					payment.IdempotencyKey,
				).WillReturnRows(sqlmock.NewRows(
					[]string{
						"id",
						"merchant_id",
						"approved",
						"status",
						"amount",
						"currency",
						"source_first_name",
						"source_last_name",
						"source_number",
						"source_bin",
						"source_card_type",
						"source_expiry_month",
						"source_expiry_year",
						"response_code",
						"response_summary",
						"reference",
						"created_at",
						"idempotency_key",
					},
				).AddRow(
					payment.Id,
					payment.MerchantID,
					payment.Approved,
					payment.Status,
					payment.Amount,
					payment.Currency,
					payment.SourceFirstName,
					payment.SourceLastName,
					payment.SourceNumber,
					payment.SourceBin,
					payment.SourceCardType,
					payment.SourceExpiryMonth,
					payment.SourceExpiryYear,
					payment.ResponseCode,
					payment.ResponseSummary,
					payment.Reference,
					payment.CreatedAt,
					payment.IdempotencyKey,
				))

				svc := NewPaymentGatewayRepository(db)
				return svc.GetByMerchantIDAndIdempotencyKey(payment.MerchantID, payment.IdempotencyKey)
			},
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}
