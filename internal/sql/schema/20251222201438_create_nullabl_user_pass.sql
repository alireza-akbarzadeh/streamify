-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Senior tip: Don't drop the column in 'Down' unless you want to lose all data. 
-- Just reverse the constraint.
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
-- +goose StatementEnd