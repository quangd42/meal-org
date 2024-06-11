-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, name, external_url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1;

-- name: UpdateRecipeByID :one
UPDATE recipes
SET name = $2, external_url = $3, updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListRecipeByUserID :many
SELECT *
FROM recipes
WHERE user_id = $1
ORDER BY updated_at DESC
LIMIT
  $2
  OFFSET $3;
