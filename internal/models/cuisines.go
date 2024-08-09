package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models/validator"
)

type CuisineRequest struct {
	Name     string     `json:"name" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

func (cr CuisineRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(cr)
}

type Cuisine struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
}
