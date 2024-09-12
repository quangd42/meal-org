-- +goose Up
ALTER TABLE recipes
ADD COLUMN external_image_url text;

-- +goose Down
ALTER TABLE recipes
DROP COLUMN external_image_url;
