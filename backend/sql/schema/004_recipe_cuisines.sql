-- +goose Up
CREATE TABLE cuisines (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) NOT NULL,
  parent_id UUID
);

CREATE TABLE recipe_cuisine (
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  cuisine_id UUID REFERENCES cuisines (id) ON DELETE CASCADE,
  recipe_id UUID REFERENCES recipes (id) ON DELETE CASCADE,
  PRIMARY KEY (recipe_id, cuisine_id)
);

-- +goose Down
DROP TABLE recipe_cuisine;
DROP TABLE cuisines;
