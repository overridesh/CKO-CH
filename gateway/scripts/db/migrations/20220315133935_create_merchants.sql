-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS merchants (
    id UUID DEFAULT gen_random_uuid(),
	name varchar(255) NOT NULL,
	apikey UUID DEFAULT gen_random_uuid(),
	active boolean DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE merchants
-- +goose StatementEnd
