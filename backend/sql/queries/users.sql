-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, username, hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
