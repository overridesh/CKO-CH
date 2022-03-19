package model

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Transaction struct {
	Id                uuid.UUID      `db:"id"`
	MerchantID        uuid.UUID      `db:"merchant_id"`
	Approved          bool           `db:"approved"`
	Status            Status         `db:"status"`
	Amount            int32          `db:"amount"`
	Currency          string         `db:"currency"`
	CreatedAt         time.Time      `db:"created_at"`
	SourceFirstName   string         `db:"source_first_name"`
	SourceLastName    string         `db:"source_last_name"`
	SourceNumber      string         `db:"source_number"`
	SourceBin         sql.NullString `db:"source_bin"`
	SourceCardType    sql.NullString `db:"source_card_type"`
	SourceExpiryMonth string         `db:"source_expiry_month"`
	SourceExpiryYear  string         `db:"source_expiry_year"`
	ResponseCode      sql.NullString `db:"response_code"`
	ResponseSummary   sql.NullString `db:"response_summary"`
	Reference         string         `db:"reference"`
	IdempotencyKey    string         `db:"idempotency_key"`
}
