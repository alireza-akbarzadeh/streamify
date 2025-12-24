-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN bio TEXT,
    ADD COLUMN phone_number VARCHAR(32) UNIQUE,
    ADD COLUMN avatar_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS avatar_url;
-- +goose StatementEnd
