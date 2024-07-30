-- +goose Up
CREATE TABLE instructions (
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  step_no INT NOT NULL,
  instruction TEXT NOT NULL,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
  PRIMARY KEY (step_no, recipe_id)
);

-- +goose Down
DROP TABLE instructions;
