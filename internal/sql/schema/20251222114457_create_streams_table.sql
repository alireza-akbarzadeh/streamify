-- +goose Up
-- +goose StatementBegin
CREATE TABLE streams (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
	device VARCHAR(50),
	started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	ended_at TIMESTAMP WITH TIME ZONE,
	completed BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS streams;
-- +goose StatementEnd