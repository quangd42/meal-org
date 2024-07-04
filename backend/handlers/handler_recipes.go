package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

	type step struct {
		StepNo      int    `json:"step_no"`
		Instruction string `json:"instruction"`
	}

	type parameters struct {
		Name              string  `json:"name"`
		ExternalURL       *string `json:"external_url"`
		Servings          int     `json:"servings"`
		Yield             *string `json:"yield"`
		CookTimeInMinutes int     `json:"cook_time_in_minutes"`
		Notes             *string `json:"notes"`
		Ingredients       []struct {
			ID       uuid.UUID `json:"id"`
			Amount   string    `json:"amount"`
			PrepNote *string   `json:"prep_note"`
		} `json:"ingredients"`
		Instructions []step `json:"instructions"`
	}

	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	// TODO: valitation of required data?

	// Create host Recipe
	recipe, err := store.Q.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		Name:              params.Name,
		ExternalUrl:       params.ExternalURL,
		Servings:          int32(params.Servings),
		Yield:             params.Yield,
		CookTimeInMinutes: int32(params.CookTimeInMinutes),
		Notes:             params.Notes,
		UserID:            userID,
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
			IngredientID: p.ID,
			Amount:       p.Amount,
			PrepNote:     p.PrepNote,
		}
	}

	// Can add as many ingredients as desired, no check here
	_, err = store.Q.AddIngredientsToRecipe(r.Context(), dbParams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Add Ingredients to host Recipe
	ingredients, err := store.Q.ListIngredientsByRecipeID(r.Context(), recipe.ID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	// Add Instructions to host Recipe
	for _, I := range params.Instructions {
		err = store.Q.AddInstructionToRecipe(r.Context(), database.AddInstructionToRecipeParams{
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			StepNo:      int32(I.StepNo),
			Instruction: I.Instruction,
			RecipeID:    recipe.ID,
		})
		if err != nil {
			respondUniqueValueError(w, err, "step_no must be unique for each step")
			return
		}
	}

	instructions, err := store.Q.ListInstructionsByRecipeID(r.Context(), recipe.ID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusCreated, createRecipeResponse(recipe, ingredients, instructions))
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

	params := &database.UpdateAllRecipeDataParams{}
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

	if err != nil {
		respondInternalServerError(w)
		return
	}

	recipe, err := database.UpdateAllRecipeData(r.Context(), store, *params)
	if err != nil {
		respondMalformedRequestError(w)
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

	recipe, err := store.Q.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		respondInternalServerError(w)
		return
	}

	if recipe.UserID != userID {
		respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	ingredients, err := store.Q.ListIngredientsByRecipeID(r.Context(), recipeID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	instructions, err := store.Q.ListInstructionsByRecipeID(r.Context(), recipe.ID)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createRecipeResponse(recipe, ingredients, instructions))
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
		respondError(w, http.StatusBadRequest, "recipe id not found")
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
