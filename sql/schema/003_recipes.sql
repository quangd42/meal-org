-- +goose Up
CREATE TABLE recipes (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  external_url TEXT,
  name TEXT NOT NULL,
  description TEXT,
  servings INT NOT NULL DEFAULT 0,
  yield TEXT,
  cook_time_in_minutes INT NOT NULL DEFAULT 0,
  notes TEXT,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE recipes;
