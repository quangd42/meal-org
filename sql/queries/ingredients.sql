-- name: CreateIngredient :one
INSERT INTO ingredients (id, created_at, updated_at, name, parent_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetIngredientByID :one
SELECT *
FROM ingredients
WHERE id = $1;

-- name: UpdateIngredientByID :one
UPDATE ingredients
SET
  name = $2,
  parent_id = $3,
  updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListIngredients :many
SELECT *
FROM ingredients
ORDER BY name;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1;
