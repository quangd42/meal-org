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

-- name: ListIngredientsByRecipeID :many
SELECT
  id,
  name,
  amount,
  instruction,
  recipe_id
FROM ingredients
JOIN recipe_ingredient ON id = ingredient_id
WHERE recipe_id = $1;

-- name: AddIngredientsToRecipe :copyfrom
INSERT INTO recipe_ingredient (
  amount, instruction, created_at, updated_at, ingredient_id, recipe_id
) VALUES ($1, $2, $3, $4, $5, $6);
