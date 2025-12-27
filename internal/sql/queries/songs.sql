-- name: CreateSong :one
INSERT INTO songs (id, title, artist_id, album_id, duration, release_date, genre, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
RETURNING *;


-- name: UpdateSong :one
UPDATE songs
SET
  title = COALESCE(sqlc.narg('title'), title),
  artist_id = COALESCE(sqlc.narg('artist_id'), artist_id),
  album_id = COALESCE(sqlc.narg('album_id'), album_id),
  duration = COALESCE(sqlc.narg('duration'), duration),
  release_date = COALESCE(sqlc.narg('release_date'), release_date),
  genre = COALESCE(sqlc.narg('genre'), genre),
  url = COALESCE(sqlc.narg('url'), url),
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetSongByID :one
SELECT * FROM songs WHERE id = $1;

-- name: ListSongs :many
SELECT * FROM songs
WHERE
  (sqlc.narg('search_title')::text IS NULL OR title ILIKE '%' || sqlc.narg('search_title') || '%')
  AND (sqlc.narg('search_artist_id')::uuid IS NULL OR artist_id = sqlc.narg('search_artist_id'))
  AND (sqlc.narg('search_genre')::text IS NULL OR genre ILIKE '%' || sqlc.narg('search_genre') || '%')
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSongs :one
SELECT COUNT(*) FROM songs
WHERE
  (sqlc.narg('search_title')::text IS NULL OR title ILIKE '%' || sqlc.narg('search_title') || '%')
  AND (sqlc.narg('search_artist_id')::uuid IS NULL OR artist_id = sqlc.narg('search_artist_id'))
  AND (sqlc.narg('search_genre')::text IS NULL OR genre ILIKE '%' || sqlc.narg('search_genre') || '%');


-- name: SoftDeleteSong :exec
UPDATE songs SET deleted_at = NOW() WHERE id = $1;

-- name: HardDeleteSong :exec
DELETE FROM songs WHERE id = $1;