-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlist_items (
	playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
	media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
	position INT,
	added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
	PRIMARY KEY (playlist_id, media_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS playlist_items;
-- +goose StatementEnd