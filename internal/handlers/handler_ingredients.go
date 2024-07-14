package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

type IngredientService interface {
	CreateIngredient(ctx context.Context, arg models.IngredientRequest) (models.Ingredient, error)
	UpdateIngredientByID(ctx context.Context, ingredientID uuid.UUID, arg models.IngredientRequest) (models.Ingredient, error)
	ListIngredients(ctx context.Context) ([]models.Ingredient, error)
	DeleteIngredient(ctx context.Context, ingredientID uuid.UUID) error
}

// TODO: restrict operation on ingredients to admin
func createIngredientHandler(is IngredientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arg, err := decodeValidate[models.IngredientRequest](r)
		if err != nil {
			respondMalformedRequestError(w)
			return
		}

		ingredient, err := is.CreateIngredient(r.Context(), arg)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusBadRequest, "invalid parent ingredient")
				return
			}
			respondDBConstraintsError(w, err, "ingredient name")
			return
		}

		respondJSON(w, http.StatusCreated, ingredient)
	}
}

func updateIngredientHandler(is IngredientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDString := chi.URLParam(r, "id")
		ingredientID, err := uuid.Parse(ingredientIDString)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid ingredient id")
			return
		}

		arg, err := decodeValidate[models.IngredientRequest](r)
		if err != nil {
			respondMalformedRequestError(w)
			return
		}

		ingredient, err := is.UpdateIngredientByID(r.Context(), ingredientID, arg)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusBadRequest, "invalid parent ingredient")
				return
			}
			respondDBConstraintsError(w, err, "ingredient name")
			return
		}

		respondJSON(w, http.StatusOK, ingredient)
	}
}

func listIngredientsHandler(is IngredientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := is.ListIngredients(r.Context())
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, ingredients)
	}
}

func deleteIngredientHandler(is IngredientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDString := chi.URLParam(r, "id")
		ingredientID, err := uuid.Parse(ingredientIDString)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid ingredient id")
			return
		}

		err = is.DeleteIngredient(r.Context(), ingredientID)
		if err != nil {
			respondDBConstraintsError(w, err, "ingredient children")
		}

		respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
	}
}
