package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &Parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	hash, err := auth.HashPassword([]byte(params.Password))
	if err != nil {
		log.Printf("error hashing password: %s\n", err)
		respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Username:  params.Username,
		Hash:      string(hash),
	})
	if err != nil {
		log.Printf("error creating new user: %s\n", err)
		respondUniqueValueError(w, err, "username")
		return
	}

	jwt, refreshToken, err := generateAndSaveAuthTokens(r, user)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusCreated, createUserResponseWithToken(user, jwt, refreshToken))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
		respondMalformedRequestError(w)
		return
	}

	hash, err := auth.HashPassword([]byte(params.Password))
	if err != nil {
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	user, err := DB.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:        userID,
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Hash:      string(hash),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondJSON(w, http.StatusOK, createUserResponse(user))
}

func forgetMeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, auth.ErrTokenNotFound.Error())
		return
	}

	err := DB.DeleteUser(r.Context(), userID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusNoContent, "user deleted")
}
