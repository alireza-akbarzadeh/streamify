-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
	status VARCHAR(20) NOT NULL,
	started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	ends_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
