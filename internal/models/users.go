package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models/validator"
)

type CreateUserRequest struct {
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password,omitempty" form:"confirm_password" validate:"omitempty,eqfield=Password"`
}

func (ur CreateUserRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(ur)
}

type UpdateUserRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

func (ur UpdateUserRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(ur)
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

func (ur LoginRequest) Validate(ctx context.Context) error {
	return validator.ValidateStruct(ur)
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

type UserWithToken struct {
	User
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (u User) WithToken(jwt, refreshToken string) UserWithToken {
	return UserWithToken{
		User:         u,
		Token:        jwt,
		RefreshToken: refreshToken,
	}
}
