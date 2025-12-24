-- +goose Up
-- +goose StatementBegin
CREATE TABLE plans (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(50) NOT NULL,
	price_cents INT NOT NULL,
	interval VARCHAR(10) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd