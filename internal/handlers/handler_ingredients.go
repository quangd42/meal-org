package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/quangd42/meal-planner/internal/database"
)

// TODO: restrict operation on ingredients to admin
func createIngredientHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string     `json:"name"`
		ParentID *uuid.UUID `json:"parent_id"`
	}
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	if params.ParentID != nil {
		_, err = store.Q.GetIngredientByID(r.Context(), *params.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent ingredient does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
	}

	ingredient, err := store.Q.CreateIngredient(r.Context(), database.CreateIngredientParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ParentID:  params.ParentID,
	})
	if err != nil {
		log.Printf("error creating new ingredient: %s\n", err)
		// Quick patch to pass tests
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code[0:2] == "23" {
			respondError(w, http.StatusBadRequest, "ingredient name")
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusCreated, createIngredientResponse(ingredient))
}

func updateIngredientHandler(w http.ResponseWriter, r *http.Request) {
	ingredientIDString := chi.URLParam(r, "id")
	ingredientID, err := uuid.Parse(ingredientIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ingredient id not found")
		return
	}

	type parameters struct {
		Name     string     `json:"name"`
		ParentID *uuid.UUID `json:"parent_id"`
	}
	params := &parameters{}
	err = json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	if params.ParentID != nil {
		_, err = store.Q.GetIngredientByID(r.Context(), *params.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent ingredient does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
	}

	ingredient, err := store.Q.UpdateIngredientByID(r.Context(), database.UpdateIngredientByIDParams{
		ID:        ingredientID,
		Name:      params.Name,
		ParentID:  params.ParentID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("error updating ingredient: %s\n", err)
		// Quick patch to pass tests
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code[0:2] == "23" {
			respondError(w, http.StatusForbidden, "cuisine id")
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createIngredientResponse(ingredient))
}

func listIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	ingredients, err := store.Q.ListIngredients(r.Context())
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "no ingredients found")
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, ingredients)
}

func deleteIngredientHandler(w http.ResponseWriter, r *http.Request) {
	ingredientIDString := chi.URLParam(r, "id")
	ingredientID, err := uuid.Parse(ingredientIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ingredient id not found")
		return
	}

	err = store.Q.DeleteIngredient(r.Context(), ingredientID)
	if err != nil {
		// Quick patch to pass tests
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code[0:2] == "23" {
			respondError(w, http.StatusForbidden, "cuisine id")
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
}