-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
	body TEXT NOT NULL,
	parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd
