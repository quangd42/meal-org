-- name: ListIngredientsByRecipeID :many
SELECT
  id,
  name,
  amount,
  prep_note,
  recipe_id
FROM ingredients
JOIN recipe_ingredient ON id = ingredient_id
WHERE recipe_id = $1;

-- name: AddIngredientsToRecipe :copyfrom
INSERT INTO recipe_ingredient (
  amount, prep_note, created_at, updated_at, ingredient_id, recipe_id
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateIngredientInRecipe :exec
UPDATE recipe_ingredient
SET
  amount = $1,
  prep_note = $2,
  updated_at = $3
WHERE
  ingredient_id = $4 AND recipe_id = $5;
