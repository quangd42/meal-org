-- name: ListIngredientsByRecipeID :many
SELECT
  id,
  name,
  amount,
  prep_note,
  recipe_id,
  index
FROM ingredients
JOIN recipe_ingredient ON id = ingredient_id
WHERE recipe_id = $1
ORDER BY index;

-- name: AddIngredientsToRecipe :copyfrom
INSERT INTO recipe_ingredient (
  amount, prep_note, created_at, updated_at, ingredient_id, recipe_id, index
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: RemoveAllIngredientsFromRecipe :exec
DELETE FROM recipe_ingredient
WHERE recipe_id = $1;
