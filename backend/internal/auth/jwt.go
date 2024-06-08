package auth

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

const (
	JWTIssuer                 = "meal_planner"
	DefaultExpirationDuration = time.Hour * 24
)

func NewJWT(id uuid.UUID, d time.Duration) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	_, ss, err := tokenAuth.Encode(map[string]interface{}{
		"Issuer":    JWTIssuer,
		"IssuedAt":  time.Now().UTC(),
		"ExpiredAt": time.Now().UTC().Add(d),
		"UserId":    id,
	})
	if err != nil {
		return "", err
	}
	return ss, nil
}
