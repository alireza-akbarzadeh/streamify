-- +goose Up
-- +goose StatementBegin

-- 1. Create enum type
CREATE TYPE user_role AS ENUM (
    'user',
    'customer',
    'admin',
    'owner'
);

-- 2. Add column with default
ALTER TABLE users
ADD COLUMN role user_role NOT NULL DEFAULT 'user';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN role;
DROP TYPE user_role;

-- +goose StatementEnd
