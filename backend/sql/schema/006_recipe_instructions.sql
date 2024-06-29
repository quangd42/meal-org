-- +goose Up
ALTER TABLE recipes
ADD COLUMN servings INT NOT NULL DEFAULT 0,
ADD COLUMN yield TEXT,
ADD COLUMN cook_time_in_minutes INT NOT NULL DEFAULT 0,
ADD COLUMN notes TEXT;

CREATE TABLE instructions (
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  step_no INT NOT NULL,
  instruction TEXT NOT NULL,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
  PRIMARY KEY (step_no, recipe_id)
);

-- +goose Down
ALTER TABLE recipes
DROP COLUMN servings,
DROP COLUMN yield,
DROP COLUMN cook_time_in_minutes,
DROP COLUMN notes;

DROP TABLE instructions;
