-- +goose Up
CREATE TABLE ingredients (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE recipe_ingredient (
  index INT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  amount TEXT NOT NULL,
  prep_note TEXT,
  ingredient_id UUID REFERENCES ingredients (id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
  PRIMARY KEY (ingredient_id, recipe_id, index)
);

-- +goose Down
DROP TABLE recipe_ingredient;
DROP TABLE ingredients;
