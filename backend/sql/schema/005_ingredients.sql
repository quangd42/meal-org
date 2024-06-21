-- +goose Up
CREATE TABLE ingredients (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) NOT NULL,
  parent_id UUID
);

CREATE TABLE recipe_ingredient (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  ingredient_id UUID REFERENCES ingredients (id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE recipe_ingredient;
DROP TABLE ingredients;
