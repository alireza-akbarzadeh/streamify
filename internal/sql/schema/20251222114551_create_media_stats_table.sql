-- +goose Up
-- +goose StatementBegin
CREATE TABLE media_stats (
	media_id UUID PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
	play_count BIGINT NOT NULL DEFAULT 0,
	like_count BIGINT NOT NULL DEFAULT 0,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS media_stats;
-- +goose StatementEnd