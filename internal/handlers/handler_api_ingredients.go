package handlers

import (
	"errors"
	"net/http"

	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

// TODO: restrict operation on ingredients to admin
func createIngredientHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arg, err := decodeJSONValidate[models.IngredientRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		ingredient, err := rs.CreateIngredient(r.Context(), arg)
		if err != nil {
			respondDBConstraintsError(w, err, "ingredient name")
			return
		}

		respondJSON(w, http.StatusCreated, ingredient)
	}
}

func updateIngredientHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientID, err := getResourceIDFromURL(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		arg, err := decodeJSONValidate[models.IngredientRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		ingredient, err := rs.UpdateIngredientByID(r.Context(), ingredientID, arg)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusBadRequest, map[string]string{"id": err.Error()})
				return
			}
			respondDBConstraintsError(w, err, "ingredient name")
			return
		}

		respondJSON(w, http.StatusOK, ingredient)
	}
}

func listIngredientsHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := rs.ListIngredients(r.Context())
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, ingredients)
	}
}

func deleteIngredientHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientID, err := getResourceIDFromURL(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = rs.DeleteIngredient(r.Context(), ingredientID)
		if err != nil {
			respondDBConstraintsError(w, err, "ingredient children")
			return
		}

		respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
	}
}
