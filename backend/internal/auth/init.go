package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

var jwtSecret string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing env settings: jwtSecret")
	}
}
