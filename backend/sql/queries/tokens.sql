-- name: SaveToken :exec
INSERT INTO tokens (
  value, created_at, expired_at, user_id
) VALUES ($1, $2, $3, $4);

-- name: GetTokenByValue :one
SELECT *
FROM tokens
WHERE value = $1;

-- name: RevokeToken :exec
UPDATE tokens
SET is_revoked = $2
WHERE value = $1;
