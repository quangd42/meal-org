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
)

var ErrRecipeNotFound = errors.New("recipe not found")

func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
		return
	}

	params := &CreateWholeRecipeParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	params.UserID = userID

	recipe, err := CreateWholeRecipe(r.Context(), store, *params)
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
		respondError(w, http.StatusBadRequest, "recipe id not found")
		return
	}

	params := &UpdateWholeRecipeParams{}
	err = json.NewDecoder(r.Body).Decode(params)
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
	params.ID = targetRecipe.ID

	err = validateRequiredParams2(params)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := UpdateWholeRecipe(r.Context(), store, *params)
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
