-- +goose Up
-- +goose StatementBegin
CREATE TABLE media (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	type VARCHAR(10) NOT NULL CHECK (type IN ('song', 'video')),
	title VARCHAR(150) NOT NULL,
	duration INT,
	artist_id UUID REFERENCES artists(id) ON DELETE SET NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS media;
-- +goose StatementEnd