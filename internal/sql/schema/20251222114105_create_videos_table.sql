-- +goose Up
-- +goose StatementBegin
CREATE TABLE videos (
	id UUID PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
	video_url TEXT NOT NULL,
	resolution VARCHAR(20),
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS videos;
-- +goose StatementEnd