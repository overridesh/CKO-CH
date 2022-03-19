-- +goose Up
-- +goose StatementBegin
INSERT INTO merchants(id, name, apikey) VALUES ('7122f444-a040-469e-9fe4-6e4f910bff93', 'TESTING MERCHANT', '19766804-2d26-4c02-ba66-751cada5cbbc');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM merchants WHERE id = '7122f444-a040-469e-9fe4-6e4f910bff93'
-- +goose StatementEnd
