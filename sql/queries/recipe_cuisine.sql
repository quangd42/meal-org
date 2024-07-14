-- name: ListCuisinesByRecipeID :many
SELECT
  id,
  name,
  recipe_id
FROM cuisines
JOIN recipe_cuisine ON id = cuisine_id
WHERE recipe_id = $1;

-- name: AddCuisinesToRecipe :exec
INSERT INTO recipe_cuisine (
  created_at, updated_at, cuisine_id, recipe_id
) VALUES ($1, $2, $3, $4);

-- name: RemoveCuisineFromRecipe :exec
DELETE FROM recipe_cuisine
WHERE
  cuisine_id = $1 AND recipe_id = $2;
