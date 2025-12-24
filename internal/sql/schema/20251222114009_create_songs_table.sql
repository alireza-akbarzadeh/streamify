-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
	id UUID PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
	album_id UUID REFERENCES albums(id) ON DELETE SET NULL,
	audio_url TEXT NOT NULL,
	bitrate INT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd