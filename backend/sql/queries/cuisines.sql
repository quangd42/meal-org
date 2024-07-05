-- name: CreateCuisine :one
INSERT INTO cuisines (id, created_at, updated_at, name, parent_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCuisineByID :one
SELECT *
FROM cuisines
WHERE id = $1;

-- name: UpdateCuisineByID :one
UPDATE cuisines
SET
  name = $2,
  parent_id = $3,
  updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListCuisines :many
SELECT *
FROM cuisines
ORDER BY name;

-- name: DeleteCuisine :exec
DELETE FROM cuisines
WHERE id = $1;
