-- +goose Up
CREATE TABLE users (
    id uuid,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name varchar(255) NOT NULL,
    username varchar(255) NOT NULL,
    password varchar(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
