-- +goose Up
ALTER TABLE recipe_ingredient
DROP CONSTRAINT recipe_ingredient_pkey;

ALTER TABLE recipe_ingredient
ADD COLUMN index INT NOT NULL;

ALTER TABLE recipe_ingredient
ADD CONSTRAINT recipe_ingredient_pkey PRIMARY KEY (
  ingredient_id, recipe_id, index
);

-- +goose Down
ALTER TABLE recipe_ingredient
DROP CONSTRAINT recipe_ingredient_pkey;

ALTER TABLE recipe_ingredient
DROP COLUMN index;

ALTER TABLE recipe_ingredient
ADD CONSTRAINT recipe_ingredient_pkey PRIMARY KEY (
  ingredient_id, recipe_id
);
