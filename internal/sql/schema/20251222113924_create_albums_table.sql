-- +goose Up
-- +goose StatementBegin
CREATE TABLE albums (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
	title VARCHAR(100) NOT NULL,
	release_date DATE,
	cover_url TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS albums;
-- +goose StatementEnd