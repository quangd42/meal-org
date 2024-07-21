package views

import (
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/components"
)

type CommonVM struct {
	Title    string
	UserID   uuid.UUID
	NavItems []components.NavItem
	Errors   map[string]any
}

func NewCommonVM(userID uuid.UUID, navItems []components.NavItem) CommonVM {
	return CommonVM{
		UserID:   userID,
		NavItems: navItems,
		Errors:   nil,
	}
}
