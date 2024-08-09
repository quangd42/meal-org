package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

type CuisineService interface {
	CreateCuisine(ctx context.Context, cr models.CuisineRequest) (models.Cuisine, error)
	UpdateCuisineByID(ctx context.Context, cuisineID uuid.UUID, cr models.CuisineRequest) (models.Cuisine, error)
	ListCuisines(ctx context.Context) ([]models.Cuisine, error)
	DeleteCuisine(ctx context.Context, cuisineID uuid.UUID) error
}

func createCuisineHandler(cs CuisineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := decodeJSONValidate[models.CuisineRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		cuisine, err := cs.CreateCuisine(r.Context(), cr)
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

func updateCuisineHandler(cs CuisineService) http.HandlerFunc {
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

		cuisine, err := cs.UpdateCuisineByID(r.Context(), cuisineID, cr)
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

func listCuisinesHandler(cs CuisineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cuisines, err := cs.ListCuisines(r.Context())
		if err != nil {
			respondInternalServerError(w)
			return
		}
		respondJSON(w, http.StatusOK, cuisines)
	}
}

func deleteCuisineHandler(cs CuisineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cuisineID, err := getResourceIDFromURL(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cs.DeleteCuisine(r.Context(), cuisineID)
		if err != nil {
			respondDBConstraintsError(w, err, "cuisine children")
			return
		}

		respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
	}
}
