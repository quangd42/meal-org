package handlers

import (
	"errors"
	"net/http"

	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

func createCuisineHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := decodeJSONValidate[models.CuisineRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		cuisine, err := rs.CreateCuisine(r.Context(), cr)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusBadRequest, map[string]string{"parent_id": "parent does not exist"})
				return
			}
			respondDBConstraintsError(w, err, "cuisine name")
			return
		}

		respondJSON(w, http.StatusCreated, cuisine)
	}
}

func updateCuisineHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cuisineID, err := getResourceIDFromURL(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		cr, err := decodeJSONValidate[models.CuisineRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		cuisine, err := rs.UpdateCuisineByID(r.Context(), cuisineID, cr)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusBadRequest, map[string]string{"parent_id": "parent does not exist"})
				return
			}
			respondDBConstraintsError(w, err, "cuisine name")
			return
		}

		respondJSON(w, http.StatusOK, cuisine)
	}
}

func listCuisinesHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cuisines, err := rs.ListCuisines(r.Context())
		if err != nil {
			respondInternalServerError(w)
			return
		}
		respondJSON(w, http.StatusOK, cuisines)
	}
}

func deleteCuisineHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cuisineID, err := getResourceIDFromURL(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = rs.DeleteCuisine(r.Context(), cuisineID)
		if err != nil {
			respondDBConstraintsError(w, err, "cuisine children")
			return
		}

		respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
	}
}
