package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/quangd42/meal-planner/internal/models"
)

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

func getPaginationParams(r *http.Request) models.RecipesPagination {
	var limit, offset int32
	limit = getPaginationParamValue(r, "limit", 20)
	offset = getPaginationParamValue(r, "offset", 0)
	return models.RecipesPagination{
		Limit:  limit,
		Offset: offset,
	}
}

func disableCacheInDevMode(next http.Handler) http.Handler {
	dev := false
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("error loading env")
	}

	devStr := os.Getenv("DEV_MODE")
	if strings.ToLower(devStr) == "true" {
		dev = true
	}

	if !dev {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	c.Render(r.Context(), w) // #nosec G104
}

func getUserIDFromCtx(ctx context.Context, sm *scs.SessionManager) (uuid.UUID, error) {
	userID, ok := sm.Get(ctx, "userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return userID, errors.New("userID could not be found")
	}
	return userID, nil
}
