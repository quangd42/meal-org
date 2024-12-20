package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-org/internal/models/validator"
)

type RecipeRequest struct {
	Name              string                `json:"name" validate:"required"`
	ExternalURL       *string               `json:"external_url" validate:"required"`
	Description       *string               `json:"description"`
	Servings          int                   `json:"servings" validate:"required"`
	Yield             *string               `json:"yield"`
	CookTimeInMinutes int                   `json:"cook_time_in_minutes" validate:"required"`
	Notes             *string               `json:"notes"`
	Cuisines          []uuid.UUID           `json:"cuisines" validate:"required,gt=0"`
	Ingredients       []IngredientInRecipe  `json:"ingredients" validate:"required,gt=0"`
	Instructions      []InstructionInRecipe `json:"instructions" validate:"required,gt=0"`
}

func (rr RecipeRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(rr)
}

type Recipe struct {
	ID                uuid.UUID             `json:"id"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	Name              string                `json:"name"`
	ExternalURL       *string               `json:"external_url"`
	ExternalImageURL  *string               `json:"external_image_url"`
	Description       *string               `json:"description"`
	UserID            uuid.UUID             `json:"user_id"`
	Servings          int                   `json:"servings"`
	Yield             *string               `json:"yield"`
	CookTimeInMinutes int                   `json:"cook_time_in_minutes"`
	Notes             *string               `json:"notes"`
	Cuisines          []CuisineInRecipe     `json:"cuisines"`
	Ingredients       []IngredientInRecipe  `json:"ingredients"`
	Instructions      []InstructionInRecipe `json:"instructions"`
}

type CuisineInRecipe struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type IngredientInRecipe struct {
	ID       uuid.UUID `json:"id"`
	Amount   string    `json:"amount"`
	PrepNote *string   `json:"prep_note"`
	Name     string    `json:"name"`
	Index    int       `json:"index"`
}

type InstructionInRecipe struct {
	StepNo      int    `json:"step_no"`
	Instruction string `json:"instruction"`
}

type RecipeInList struct {
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Name              string    `json:"name"`
	ExternalURL       *string   `json:"external_url"`
	ExternalImageURL  *string   `json:"external_image_url"`
	Description       *string   `json:"description"`
	Cuisines          string    `json:"cuisines"`
	UserID            uuid.UUID `json:"user_id"`
	Servings          int       `json:"servings"`
	Yield             *string   `json:"yield"`
	CookTimeInMinutes int       `json:"cook_time_in_minutes"`
}

type RecipesPagination struct {
	Limit  int32
	Offset int32
}
