// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: recipes.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRecipe = `-- name: CreateRecipe :one
INSERT INTO recipes (
  id,
  created_at,
  updated_at,
  name,
  external_url,
  user_id,
  servings,
  yield,
  cook_time_in_minutes,
  notes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at, updated_at, external_url, name, user_id, servings, yield, cook_time_in_minutes, notes
`

type CreateRecipeParams struct {
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Name              string    `json:"name"`
	ExternalUrl       *string   `json:"external_url"`
	UserID            uuid.UUID `json:"user_id"`
	Servings          int32     `json:"servings"`
	Yield             *string   `json:"yield"`
	CookTimeInMinutes int32     `json:"cook_time_in_minutes"`
	Notes             *string   `json:"notes"`
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error) {
	row := q.db.QueryRow(ctx, createRecipe,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.ExternalUrl,
		arg.UserID,
		arg.Servings,
		arg.Yield,
		arg.CookTimeInMinutes,
		arg.Notes,
	)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExternalUrl,
		&i.Name,
		&i.UserID,
		&i.Servings,
		&i.Yield,
		&i.CookTimeInMinutes,
		&i.Notes,
	)
	return i, err
}

const deleteRecipe = `-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1
`

func (q *Queries) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRecipe, id)
	return err
}

const getRecipeByID = `-- name: GetRecipeByID :one
SELECT id, created_at, updated_at, external_url, name, user_id, servings, yield, cook_time_in_minutes, notes FROM recipes
WHERE id = $1
`

func (q *Queries) GetRecipeByID(ctx context.Context, id uuid.UUID) (Recipe, error) {
	row := q.db.QueryRow(ctx, getRecipeByID, id)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExternalUrl,
		&i.Name,
		&i.UserID,
		&i.Servings,
		&i.Yield,
		&i.CookTimeInMinutes,
		&i.Notes,
	)
	return i, err
}

const listRecipesByUserID = `-- name: ListRecipesByUserID :many
SELECT id, created_at, updated_at, external_url, name, user_id, servings, yield, cook_time_in_minutes, notes
FROM recipes
WHERE user_id = $1
ORDER BY name
LIMIT
  $2
  OFFSET $3
`

type ListRecipesByUserIDParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListRecipesByUserID(ctx context.Context, arg ListRecipesByUserIDParams) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, listRecipesByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExternalUrl,
			&i.Name,
			&i.UserID,
			&i.Servings,
			&i.Yield,
			&i.CookTimeInMinutes,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRecipeByID = `-- name: UpdateRecipeByID :one
UPDATE recipes
SET
  name = $2,
  external_url = $3,
  updated_at = $4,
  servings = $5,
  yield = $6,
  cook_time_in_minutes = $7,
  notes = $8
WHERE id = $1
RETURNING id, created_at, updated_at, external_url, name, user_id, servings, yield, cook_time_in_minutes, notes
`

type UpdateRecipeByIDParams struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	ExternalUrl       *string   `json:"external_url"`
	UpdatedAt         time.Time `json:"updated_at"`
	Servings          int32     `json:"servings"`
	Yield             *string   `json:"yield"`
	CookTimeInMinutes int32     `json:"cook_time_in_minutes"`
	Notes             *string   `json:"notes"`
}

func (q *Queries) UpdateRecipeByID(ctx context.Context, arg UpdateRecipeByIDParams) (Recipe, error) {
	row := q.db.QueryRow(ctx, updateRecipeByID,
		arg.ID,
		arg.Name,
		arg.ExternalUrl,
		arg.UpdatedAt,
		arg.Servings,
		arg.Yield,
		arg.CookTimeInMinutes,
		arg.Notes,
	)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExternalUrl,
		&i.Name,
		&i.UserID,
		&i.Servings,
		&i.Yield,
		&i.CookTimeInMinutes,
		&i.Notes,
	)
	return i, err
}
