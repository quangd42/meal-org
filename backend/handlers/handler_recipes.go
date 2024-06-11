package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
)

var ErrRecipeNotFound = errors.New("recipe not found")

func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	type parameters struct {
		Name        string `json:"name"`
		ExternalURL string `json:"external_url,omitempty"`
	}

	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	recipe, err := DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		ExternalUrl: params.ExternalURL,
		UserID:      userID,
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, recipe)
}

func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	type parameters struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		ExternalURL string    `json:"external_url,omitempty"`
	}

	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	targetRecipe, err := DB.GetRecipeByID(r.Context(), params.ID)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrRecipeNotFound.Error())
		return
	}
	if targetRecipe.UserID != userID {
		respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	recipe, err := DB.UpdateRecipeByID(r.Context(), database.UpdateRecipeByIDParams{
		ID:          targetRecipe.ID,
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		ExternalUrl: params.ExternalURL,
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, recipe)
}

func listRecipesHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	var limit, offset int32
	limitStr := chi.URLParam(r, "limit")
	limit64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		limit = 20
	} else {
		limit = int32(limit64)
	}
	offsetStr := chi.URLParam(r, "offset")
	offset64, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		offset = 0
	} else {
		offset = int32(offset64)
	}

	recipes, err := DB.ListRecipeByUserID(r.Context(), database.ListRecipeByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, recipes)
}
