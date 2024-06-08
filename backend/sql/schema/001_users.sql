-- +goose Up
CREATE TABLE users (
    id uuid UNIQUE NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name varchar(255) NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    hash varchar(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
