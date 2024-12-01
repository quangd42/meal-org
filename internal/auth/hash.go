package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	JWTIssuer                 = "meal_planner"
	ExpirationDurationAccess  = time.Hour
	ExpirationDurationDefault = time.Hour * 24
	ExpirationDurationRefresh = time.Hour * 24 * 60
)

var (
	ErrTokenNotFound    = errors.New("token not found")
	ErrTokenInvalid     = errors.New("token expired or invalidated")
	ErrClaimTypeInvalid = errors.New("claim type cannot be verified")
)

func HashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func ValidateHash(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
