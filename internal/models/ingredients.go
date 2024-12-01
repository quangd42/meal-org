package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-org/internal/models/validator"
)

type IngredientRequest struct {
	Name string `json:"name" validate:"required"`
}

func (ir IngredientRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(ir)
}

type Ingredient struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}
