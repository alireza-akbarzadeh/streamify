-- name: UpdateUserProfile :exec
UPDATE users
SET
  first_name = COALESCE(sqlc.narg('first_name'), first_name),
  last_name = COALESCE(sqlc.narg('last_name'), last_name),
  bio = COALESCE(sqlc.narg('bio'), bio),
  avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
  phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
  updated_at = NOW()
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
  username, email, password_hash, verification_token, verification_expires_at, first_name, last_name, bio, phone_number, avatar_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2 WHERE id = $1;

UPDATE users
SET first_name = $2,
    last_name = $3,
    bio = $4,
    avatar_url = $5,
    username = $6,
    phone_number = $7
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: VerifyUserByToken :one
UPDATE users
SET is_verified = TRUE,
    verification_token = NULL,
    verification_expires_at = NULL
WHERE verification_token = $1
  AND verification_expires_at > NOW()
RETURNING id, email, username;

-- name: GetUserByVerificationToken :one
SELECT * FROM users WHERE verification_token = $1 LIMIT 1;

-- name: VerifyUserByTokenByID :exec
UPDATE users
SET is_verified = TRUE,
    verification_token = NULL,
    verification_expires_at = NULL
WHERE id = $1;

-- name: LockUser :exec
UPDATE users SET is_locked = TRUE WHERE id = $1;

-- name: UnlockUser :exec
UPDATE users SET is_locked = FALSE WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE
  (sqlc.narg('search_username')::text IS NULL OR username ILIKE '%' || sqlc.narg('search_username') || '%')
  AND (sqlc.narg('search_email')::text IS NULL OR email ILIKE '%' || sqlc.narg('search_email') || '%')
  AND (sqlc.narg('search_phone_number')::text IS NULL OR phone_number ILIKE '%' || sqlc.narg('search_phone_number') || '%')
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE
  (sqlc.narg('search_username')::text IS NULL OR username ILIKE '%' || sqlc.narg('search_username') || '%')
  AND (sqlc.narg('search_email')::text IS NULL OR email ILIKE '%' || sqlc.narg('search_email') || '%')
  AND (sqlc.narg('search_phone_number')::text IS NULL OR phone_number ILIKE '%' || sqlc.narg('search_phone_number') || '%');

