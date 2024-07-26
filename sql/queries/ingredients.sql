-- name: CreateIngredient :one
INSERT INTO ingredients (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetIngredientByID :one
SELECT *
FROM ingredients
WHERE id = $1;

-- name: UpdateIngredientByID :one
UPDATE ingredients
SET
  name = $2,
  updated_at = $3
WHERE id = $1
RETURNING *;

-- name: ListIngredients :many
SELECT *
FROM ingredients
ORDER BY name;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = $1;
