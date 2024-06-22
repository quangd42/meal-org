package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		Ingredients []struct {
			ID          uuid.UUID `json:"id"`
			Amount      string    `json:"amount"`
			Instruction string    `json:"instruction"`
		} `json:"ingredients"`
	}

	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	recipe, err := DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:          NewUUID(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		ExternalUrl: params.ExternalURL,
		UserID:      pgUUID(userID),
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	dbParams := make([]database.AddIngredientsToRecipeParams, len(params.Ingredients))
	for i, p := range params.Ingredients {
		dbParams[i] = database.AddIngredientsToRecipeParams{
			RecipeID:     recipe.ID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			IngredientID: pgUUID(p.ID),
			Amount:       p.Amount,
			Instruction:  pgtype.Text{String: p.Instruction, Valid: p.Instruction != ""},
		}
	}

	// Can add as many ingredients as desired, no check here
	_, err = DB.AddIngredientsToRecipe(r.Context(), dbParams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ingredients, err := DB.ListIngredientsByRecipeID(r.Context(), recipe.ID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createRecipeResponse(recipe, ingredients))
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

	type parameters struct {
		Name        string `json:"name"`
		ExternalURL string `json:"external_url,omitempty"`
		Ingredients []IngredientInRecipe
	}

	params := &parameters{}
	err = json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	targetRecipe, err := DB.GetRecipeByID(r.Context(), pgUUID(recipeID))
	if err != nil {
		respondError(w, http.StatusNotFound, ErrRecipeNotFound.Error())
		return
	}
	if targetRecipe.UserID.Bytes != userID {
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

	ingredients, err := DB.ListIngredientsByRecipeID(r.Context(), pgUUID(recipeID))
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createRecipeResponse(recipe, ingredients))
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
	recipes, err := DB.ListRecipesByUserID(r.Context(), database.ListRecipesByUserIDParams{
		UserID: pgUUID(userID),
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

	recipe, err := DB.GetRecipeByID(r.Context(), pgUUID(recipeID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		respondInternalServerError(w)
		return
	}

	if recipe.UserID.Bytes != userID {
		respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	ingredients, err := DB.ListIngredientsByRecipeID(r.Context(), pgUUID(recipeID))
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createRecipeResponse(recipe, ingredients))
}
