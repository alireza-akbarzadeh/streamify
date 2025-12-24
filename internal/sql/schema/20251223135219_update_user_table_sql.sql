-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    RENAME COLUMN name TO username;

ALTER TABLE users
    ADD COLUMN first_name VARCHAR(100),
    ADD COLUMN last_name VARCHAR(100),
    ADD COLUMN is_locked BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
ALTER TABLE users DROP COLUMN IF EXISTS is_locked;
ALTER TABLE users RENAME COLUMN username TO name;
-- +goose StatementEnd