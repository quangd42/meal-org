package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-org/internal/auth"
	"github.com/quangd42/meal-org/internal/database"
	"github.com/quangd42/meal-org/internal/models"
)

var ErrHashPassword = errors.New("error hashing password")

type UserService struct {
	store *database.Store
}

func NewUserService(store *database.Store) UserService {
	return UserService{store: store}
}

func (us UserService) CreateUser(ctx context.Context, ur models.CreateUserRequest) (models.User, error) {
	var u models.User
	hash, err := auth.HashPassword([]byte(ur.Password))
	if err != nil {
		log.Printf("error hashing password: %s\n", err)
		return u, ErrHashPassword
	}

	user, err := us.store.Q.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     ur.Email,
		Hash:      string(hash),
	})
	if err != nil {
		log.Printf("error creating new user: %s\n", err)
		return u, ErrDBConstraint
	}

	u = genUserResponse(user)

	return u, nil
}

func (us UserService) UpdateUserByID(ctx context.Context, userID uuid.UUID, ur models.UpdateUserRequest) (models.User, error) {
	var u models.User

	hash, err := auth.HashPassword([]byte(ur.Password))
	if err != nil {
		log.Printf("error hashing password: %s\n", err)
		return u, ErrHashPassword
	}

	user, err := us.store.Q.UpdateUserByID(ctx, database.UpdateUserByIDParams{
		ID:        userID,
		UpdatedAt: time.Now().UTC(),
		Hash:      string(hash),
	})
	if err != nil {
		return u, ErrDBConstraint
	}

	u = genUserResponse(user)

	return u, nil
}

func (us UserService) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	err := us.store.Q.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func genUserResponse(u database.User) models.User {
	return models.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		Email:     u.Email,
	}
}
