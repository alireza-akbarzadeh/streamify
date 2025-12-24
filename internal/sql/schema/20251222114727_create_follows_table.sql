-- +goose Up
-- +goose StatementBegin
CREATE TABLE follows (
	follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	PRIMARY KEY (follower_id, artist_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS follows;
-- +goose StatementEnd