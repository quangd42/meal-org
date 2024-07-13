package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/middleware"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
	"github.com/quangd42/meal-planner/internal/services/auth"
)

type RecipeService interface {
	CreateRecipe(ctx context.Context, userID uuid.UUID, rr models.RecipeRequest) (models.Recipe, error)
	UpdateRecipeByID(ctx context.Context, userID, recipeID uuid.UUID, rr models.RecipeRequest) (models.Recipe, error)
	GetRecipeByID(ctx context.Context, recipeID uuid.UUID) (models.Recipe, error)
	ListRecipesByUserID(ctx context.Context, userID uuid.UUID, pgn models.RecipesPagination) ([]models.RecipeInList, error)
	DeleteRecipeByID(ctx context.Context, recipeID uuid.UUID) error
}

// TODO: allow for uploading images
func createRecipeHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
		if !ok {
			respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
			return
		}

		rr, err := decodeValidate[models.RecipeRequest](r)
		if err != nil {
			respondMalformedRequestError(w)
			return
		}

		recipe, err := rs.CreateRecipe(r.Context(), userID, rr)
		if err != nil {
			respondDBConstraintsError(w, err, "ingredient_id, step_no")
			return
		}

		respondJSON(w, http.StatusCreated, recipe)
	}
}

func updateRecipeHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		rr, err := decodeValidate[models.RecipeRequest](r)
		if err != nil {
			respondMalformedRequestError(w)
			return
		}

		recipe, err := rs.UpdateRecipeByID(r.Context(), userID, recipeID, rr)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusNotFound, services.ErrResourceNotFound.Error())
				return
			}
			if errors.Is(err, services.ErrUnauthorized) {
				respondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			respondDBConstraintsError(w, err, "cuisine_id, step_no")
			return
		}

		respondJSON(w, http.StatusOK, recipe)
	}
}

func listRecipesHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDCtxKey).(uuid.UUID)
		if !ok {
			respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
			return
		}
		pgn := getPaginationParams(r)
		recipes, err := rs.ListRecipesByUserID(r.Context(), userID, pgn)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, recipes)
	}
}

func getRecipeHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		recipe, err := rs.GetRecipeByID(r.Context(), recipeID)
		if err != nil {
			if errors.Is(err, services.ErrResourceNotFound) {
				respondError(w, http.StatusNotFound, services.ErrResourceNotFound.Error())
				return
			}
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, recipe)
	}
}

// TODO: unit testing delete Recipe:
// make sure that instructions and ingredient links are deleted
func deleteRecipeHandler(rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = rs.DeleteRecipeByID(r.Context(), recipeID)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
	}
}
