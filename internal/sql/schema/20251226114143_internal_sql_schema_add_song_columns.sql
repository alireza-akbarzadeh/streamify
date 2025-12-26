-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs
    ADD COLUMN title TEXT NOT NULL DEFAULT '',
    ADD COLUMN artist_id UUID REFERENCES artists(id) ON DELETE SET NULL,
    ADD COLUMN duration INT,
    ADD COLUMN release_date DATE,
    ADD COLUMN genre TEXT,
    ADD COLUMN url TEXT,
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs
    DROP COLUMN title,
    DROP COLUMN artist_id,
    DROP COLUMN duration,
    DROP COLUMN release_date,
    DROP COLUMN genre,
    DROP COLUMN url,
    DROP COLUMN updated_at,
    DROP COLUMN deleted_at;
-- +goose StatementEnd

