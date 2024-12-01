package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/quangd42/meal-org/internal/auth"
	"github.com/quangd42/meal-org/internal/database"
	"github.com/quangd42/meal-org/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	jwtSecret string
	store     *database.Store
}

func NewAuthService(store *database.Store, jwtSecret string) Auth {
	return Auth{
		jwtSecret: jwtSecret,
		store:     store,
	}
}

func (as Auth) GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	jwt, err := auth.CreateJWT(as.jwtSecret, userID, auth.ExpirationDurationAccess)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (as Auth) GenerateAndSaveRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = as.store.Q.SaveToken(ctx, database.SaveTokenParams{
		Value:     refreshToken,
		CreatedAt: time.Now().UTC(),
		ExpiredAt: time.Now().Add(auth.ExpirationDurationRefresh),
		UserID:    userID,
	})
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (as Auth) ValidateRefreshToken(ctx context.Context, refreshToken string) (userID uuid.UUID, err error) {
	token, err := as.store.Q.GetTokenByValue(ctx, refreshToken)
	if err != nil {
		return userID, err
	}

	if token.ExpiredAt.Before(time.Now().UTC()) {
		token.IsRevoked = true
		err = as.store.Q.RevokeToken(ctx, database.RevokeTokenParams{
			Value:     refreshToken,
			IsRevoked: true,
		})
		if err != nil {
			return userID, err
		}
	}
	if token.IsRevoked {
		return userID, auth.ErrTokenInvalid
	}

	return token.UserID, nil
}

func (as Auth) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	_, err := as.store.Q.GetTokenByValue(ctx, refreshToken)
	if err != nil {
		return nil
	}

	err = as.store.Q.RevokeToken(ctx, database.RevokeTokenParams{
		Value:     refreshToken,
		IsRevoked: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (as Auth) Login(ctx context.Context, lr models.LoginRequest) (models.User, error) {
	var u models.User
	user, err := as.store.Q.GetUserByEmail(ctx, lr.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, pgx.ErrNoRows
		}
		return u, err
	}

	err = auth.ValidateHash([]byte(user.Hash), []byte(lr.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return u, bcrypt.ErrMismatchedHashAndPassword
		}
		return u, err
	}

	u = genUserResponse(user)

	return u, nil
}
