-- +goose Up
CREATE TABLE users (
  id UUID UNIQUE PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  hash VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
