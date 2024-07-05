package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
	"github.com/quangd42/meal-planner/backend/internal/models"
)

var ErrRecipeNotFound = errors.New("recipe not found")

// TODO: allow for uploading images
func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	rr := &models.RecipeRequest{}
	err := json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	err = rr.Validate()
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	params := CreateWholeRecipeParams{
		UserID:        userID,
		RecipeRequest: *rr,
	}

	recipe, err := CreateWholeRecipe(r.Context(), store, params)
	if err != nil {
		respondDBConstraintsError(w, err, "ingredient_id, step_no")
		return
	}

	respondJSON(w, http.StatusCreated, recipe)
}

func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	recipeIDString := chi.URLParam(r, "id")
	recipeID, err := uuid.Parse(recipeIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	rr := &models.RecipeRequest{}
	err = json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	err = rr.Validate()
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	// Check if the recipe belongs to the user
	// NOTE: this perhaps should be in a middleware
	targetRecipe, err := store.Q.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrRecipeNotFound.Error())
		return
	}
	if targetRecipe.UserID != userID {
		respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	params := UpdateWholeRecipeParams{
		ID:            recipeID,
		RecipeRequest: *rr,
	}

	recipe, err := UpdateWholeRecipe(r.Context(), store, params)
	if err != nil {
		respondDBConstraintsError(w, err, "ingredient_id, step_no")
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
	limit = getPaginationParamValue(r, "limit", 20)
	offset = getPaginationParamValue(r, "offset", 0)
	recipes, err := store.Q.ListRecipesByUserID(r.Context(), database.ListRecipesByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, recipes)
}

func getRecipeHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	recipeIDString := chi.URLParam(r, "id")
	recipeID, err := uuid.Parse(recipeIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	recipe, err := GetWholeRecipe(r.Context(), store, recipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, ErrRecipeNotFound.Error())
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, recipe)
}

// TODO: unit testing delete Recipe: make sure that instructions and ingredient links
// are deleted
func deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	recipeIDString := chi.URLParam(r, "id")
	recipeID, err := uuid.Parse(recipeIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid recipe id")
		return
	}

	err = store.Q.DeleteRecipe(r.Context(), database.DeleteRecipeParams{
		UserID: userID,
		ID:     recipeID,
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
