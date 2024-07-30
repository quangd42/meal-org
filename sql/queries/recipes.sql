-- name: CreateRecipe :one
INSERT INTO recipes (
  id,
  created_at,
  updated_at,
  name,
  description,
  external_url,
  user_id,
  servings,
  yield,
  cook_time_in_minutes,
  notes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1;

-- name: UpdateRecipeByID :one
UPDATE recipes
SET
  name = $2,
  external_url = $3,
  description = $9,
  updated_at = $4,
  servings = $5,
  yield = $6,
  cook_time_in_minutes = $7,
  notes = $8
WHERE id = $1
RETURNING *;

-- name: ListRecipesByUserID :many
SELECT *
FROM recipes
WHERE user_id = $1
ORDER BY name
LIMIT
  $2
  OFFSET $3;

-- name: ListRecipesWithCuisinesByUserID :many
SELECT
  r.*,
  string_agg(c.name, ', ') AS cuisines
FROM
  recipes r
LEFT JOIN
  recipe_cuisine rc ON r.id = rc.recipe_id
LEFT JOIN
  cuisines c ON rc.cuisine_id = c.id
WHERE
  r.user_id = $1
GROUP BY
  r.id
LIMIT
  $2
  OFFSET $3;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;
