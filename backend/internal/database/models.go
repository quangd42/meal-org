// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Cuisine struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ParentID  uuid.UUID `json:"parent_id"`
}

type Ingredient struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ParentID  uuid.UUID `json:"parent_id"`
}

type Recipe struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExternalUrl string    `json:"external_url"`
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
}

type RecipeCuisine struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CuisineID uuid.UUID `json:"cuisine_id"`
	RecipeID  uuid.UUID `json:"recipe_id"`
}

type RecipeIngredient struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IngredientID uuid.UUID `json:"ingredient_id"`
	RecipeID     uuid.UUID `json:"recipe_id"`
}

type Token struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	IsRevoked bool      `json:"is_revoked"`
	UserID    uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Hash      string    `json:"hash"`
}
