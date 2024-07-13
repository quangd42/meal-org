package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ur CreateUserRequest) Validate(ctx context.Context) error {
	return validate.Struct(ur)
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ur UpdateUserRequest) Validate(ctx context.Context) error {
	return validate.Struct(ur)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ur LoginRequest) Validate(ctx context.Context) error {
	return validate.Struct(ur)
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
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
