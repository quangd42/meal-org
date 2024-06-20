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
SET name = $2,
    parent_id = $3,
    updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListIngredientsByRecipeID :many
SELECT *
FROM ingredients
WHERE id in (
        SELECT ingredient_id
        FROM recipe_ingredient
        WHERE recipe_id = $1
    )
ORDER BY ingredients.name;

-- name: AddIngredientsToRecipe :batchmany
-- TODO: add this
