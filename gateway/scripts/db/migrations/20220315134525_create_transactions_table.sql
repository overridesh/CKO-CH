-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id UUID DEFAULT gen_random_uuid(),
	merchant_id UUID NOT NULL,
    approved BOOLEAN DEFAULT FALSE,
    status VARCHAR(255) DEFAULT 'pending',
    amount integer NOT NULL, 
    currency VARCHAR(3) NOT NULL,
    source_first_name VARCHAR(255) NOT NULL,
    source_last_name VARCHAR(255) NOT NULL,
    source_number VARCHAR(19) NOT NULL,
    source_bin VARCHAR(6),
    source_card_type VARCHAR(255),
    source_expiry_month VARCHAR(2) NOT NULL,
    source_expiry_year VARCHAR(4) NOT NULL,
    idempotency_key VARCHAR(255),
    response_code VARCHAR(255),
    response_summary VARCHAR(255),
    reference VARCHAR(255),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
	CONSTRAINT FK_merchant_Id FOREIGN KEY (merchant_id)
    REFERENCES merchants(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions
-- +goose StatementEnd
