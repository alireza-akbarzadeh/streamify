-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS songs CASCADE;

CREATE TABLE songs (
    id UUID PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
    album_id UUID REFERENCES albums(id) ON DELETE SET NULL,
    audio_url TEXT NOT NULL,
    bitrate INT,
    title TEXT NOT NULL DEFAULT '',
    artist_id UUID REFERENCES artists(id) ON DELETE SET NULL,
    duration INT,
    release_date DATE,
    genre TEXT,
    url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd