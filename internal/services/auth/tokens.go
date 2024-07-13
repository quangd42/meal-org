package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uuid.UUID `json:"userID"`
	jwt.RegisteredClaims
}

func (uc *UserClaims) GetUserID() uuid.UUID {
	return uc.UserID
}

func CreateJWT(userID uuid.UUID, d time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := UserClaims{
		userID,
		jwt.RegisteredClaims{
			Issuer:    JWTIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(d)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func VerifyJWT(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return uuid.UUID{}, ErrClaimTypeInvalid
	}

	return claims.UserID, nil
}

func GetHeaderToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if len(header) < 7 || strings.ToLower(header[0:6]) != "bearer" {
		return "", ErrTokenNotFound
	}
	return header[7:], nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	tokenString := base64.StdEncoding.EncodeToString(b)
	return tokenString, nil
}
