package views

import (
	"github.com/google/uuid"
)

type CommonVM struct {
	Title    string
	UserID   uuid.UUID
	NavItems []NavItem
	Errors   map[string]any
}

func NewCommonVM(userID uuid.UUID, navItems []NavItem) CommonVM {
	return CommonVM{
		UserID:   userID,
		NavItems: navItems,
		Errors:   nil,
	}
}
