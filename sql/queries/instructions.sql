-- name: AddInstructionToRecipe :exec
INSERT INTO instructions (
  created_at, updated_at, instruction, step_no, recipe_id
) VALUES ($1, $2, $3, $4, $5);

-- name: ListInstructionsByRecipeID :many
SELECT *
FROM instructions
WHERE recipe_id = $1
ORDER BY step_no ASC;

-- name: UpdateInstructionByID :exec
UPDATE instructions
SET updated_at = $3, instruction = $4
WHERE step_no = $1 AND recipe_id = $2;

-- name: DeleteInstructionByID :exec
DELETE FROM instructions
WHERE step_no = $1 AND recipe_id = $2;
