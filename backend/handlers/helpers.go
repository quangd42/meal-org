package handlers

import (
	"net/http"
	"time"

	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

func GenerateAndSaveAuthTokens(r *http.Request, user database.User) (string, string, error) {
	jwt, err := auth.CreateJWT(user, auth.ExpirationDurationAccess)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = DB.SaveToken(r.Context(), database.SaveTokenParams{
		Value:     refreshToken,
		CreatedAt: time.Now().UTC(),
		ExpiredAt: time.Now().Add(auth.ExpirationDurationRefresh),
		UserID:    user.ID,
	})
	if err != nil {
		return "", "", err
	}

	return jwt, refreshToken, nil
}
