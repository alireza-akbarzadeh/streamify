-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_profiles (
	user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	avatar_url TEXT,
	bio TEXT,
	country VARCHAR(50),
	birth_date DATE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
