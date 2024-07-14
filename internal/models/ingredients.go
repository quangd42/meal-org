package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IngredientRequest struct {
	Name     string     `json:"name" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

func (ingR IngredientRequest) Validate(ctx context.Context) error {
	return validate.Struct(ingR)
}

type Ingredient struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
}
