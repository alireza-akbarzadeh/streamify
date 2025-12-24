-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
-- +goose StatementEnd
