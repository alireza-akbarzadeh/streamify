-- +goose Up
-- +goose StatementBegin

-- 1. Create new enum
CREATE TYPE user_role_new AS ENUM (
    'user',
    'customer',
    'admin',
    'owner'
);

-- 2. Update column to new enum
ALTER TABLE users
ALTER COLUMN role DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN role TYPE user_role_new
USING (
    CASE role
        WHEN 'regular' THEN 'user'::user_role_new
        ELSE role::text::user_role_new
    END
);

-- 3. Set new default
ALTER TABLE users
ALTER COLUMN role SET DEFAULT 'user';

-- 4. Cleanup
DROP TYPE user_role;
ALTER TYPE user_role_new RENAME TO user_role;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Reverse back (optional but correct)

CREATE TYPE user_role_old AS ENUM (
    'regular',
    'customer',
    'admin',
    'owner'
);

ALTER TABLE users
ALTER COLUMN role DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN role TYPE user_role_old
USING (
    CASE role
        WHEN 'user' THEN 'regular'::user_role_old
        ELSE role::text::user_role_old
    END
);

ALTER TABLE users
ALTER COLUMN role SET DEFAULT 'regular';

DROP TYPE user_role;
ALTER TYPE user_role_old RENAME TO user_role;

-- +goose StatementEnd
