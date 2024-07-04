package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                uuid.UUID             `json:"id"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	Name              string                `json:"name"`
	ExternalUrl       *string               `json:"external_url"`
	UserID            uuid.UUID             `json:"user_id"`
	Servings          int                   `json:"servings"`
	Yield             *string               `json:"yield"`
	CookTimeInMinutes int                   `json:"cook_time_in_minutes"`
	Notes             *string               `json:"notes"`
	Ingredients       []IngredientInRecipe  `json:"ingredients"`
	Instructions      []InstructionInRecipe `json:"instructions"`
}

type IngredientInRecipe struct {
	ID       uuid.UUID `json:"id"`
	Amount   string    `json:"amount"`
	PrepNote *string   `json:"prep_note"`
	Name     string    `json:"name"`
}

type InstructionInRecipe struct {
	StepNo      int    `json:"step_no"`
	Instruction string `json:"instruction"`
}
