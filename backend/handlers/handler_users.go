package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
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
			log.Printf("error hashing password: %s\n", err)
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
			log.Printf("error creating new user: %s\n", err)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23505") {
				respondError(w, http.StatusBadRequest, "User already exists")
				return
			}
			respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		token, err := auth.CreateJWT(user, auth.DefaultExpirationDuration)
		if err != nil {
			log.Printf("error creating new JWT: %s\n", err)
			respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondJSON(w, http.StatusOK, createResponseUserWithToken(user, token))
	}
}

func UpdateUserHandler(c *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
		if !ok {
			respondError(w, http.StatusUnauthorized, auth.ErrTokenNotFound.Error())
			return
		}

		decoder := json.NewDecoder(r.Body)
		type Parameters struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		params := &Parameters{}
		err := decoder.Decode(params)
		if err != nil {
			log.Printf("error decoding: %s\n", err.Error())
			respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		hash, err := auth.HashPassword([]byte(params.Password))
		if err != nil {
			respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		user, err := c.DB.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
			ID:        userID,
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Hash:      string(hash),
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondJSON(w, http.StatusOK, createResponseUser(user))
	}
}
