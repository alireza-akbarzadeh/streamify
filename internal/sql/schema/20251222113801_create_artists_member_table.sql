-- +goose Up
-- +goose StatementBegin
CREATE TABLE artist_members (
	artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role VARCHAR(50),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	PRIMARY KEY (artist_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS artist_members;
-- +goose StatementEnd