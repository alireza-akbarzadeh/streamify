-- +goose Up
-- +goose StatementBegin
-- 1. Add column as nullable first to avoid breaking existing rows
ALTER TABLE users ADD COLUMN password TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS password;
-- +goose StatementEnd
