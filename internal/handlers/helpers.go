package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/services/auth"
)

func generateAndSaveAuthTokens(r *http.Request, user database.User) (string, string, error) {
	jwt, err := auth.CreateJWT(user, auth.ExpirationDurationAccess)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = store.Q.SaveToken(r.Context(), database.SaveTokenParams{
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

func getPaginationParamValue(r *http.Request, name string, defaultValue int32) int32 {
	val := int32(defaultValue)
	paramStr := r.URL.Query().Get(name)
	if paramStr == "" {
		return val
	}
	param64, err := strconv.ParseInt(paramStr, 10, 32)
	if err != nil {
		return val
	}
	val = int32(param64)
	return val
}
