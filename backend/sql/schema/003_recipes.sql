-- +goose Up
CREATE TABLE recipes (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  external_url TEXT NOT NULL,
  name TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE recipes;
