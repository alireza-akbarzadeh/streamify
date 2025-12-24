-- +goose Up
-- +goose StatementBegin
CREATE TABLE likes (
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, media_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS likes;
-- +goose StatementEnd