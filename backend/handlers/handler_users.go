package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

func CreateUserHandler(c *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			Name     string `json:"name"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		params := &Parameters{}
		err := decoder.Decode(params)
		if err != nil {
			respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		hash, err := auth.HashPassword([]byte(params.Password))
		if err != nil {
			respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Username:  params.Username,
			Hash:      string(hash),
		})
		if err != nil {
			respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		respondJSON(w, http.StatusOK, createResponseUser(user))
	}
}
