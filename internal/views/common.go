package views

import (
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models"
)

type CommonVM struct {
	Title    string
	UserID   uuid.UUID
	NavItems []models.NavItem
	Errors   map[string]any
}

func NewCommonVM(userID uuid.UUID, navItems []models.NavItem) CommonVM {
	return CommonVM{
		UserID:   userID,
		NavItems: navItems,
		Errors:   nil,
	}
}
