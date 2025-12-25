-- name: CreateSession :one
INSERT INTO user_sessions (
    user_id, refresh_token, ip_address, user_agent, expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM user_sessions 
WHERE refresh_token = $1 LIMIT 1;

-- name: DeleteSessionByToken :exec
DELETE FROM user_sessions WHERE refresh_token = $1;

-- name: DeleteAllUserSessions :exec
DELETE FROM user_sessions WHERE user_id = $1;

-- name: GetSessionByID :one
SELECT * FROM user_sessions WHERE id = $1 LIMIT 1;

-- name: DeleteSessionByID :exec
DELETE FROM user_sessions WHERE id = $1;


