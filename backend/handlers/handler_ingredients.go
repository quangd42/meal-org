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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

// TODO: only allow admin user
func createIngredientHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string    `json:"name"`
		ParentID uuid.UUID `json:"parent_id"`
	}
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	ingredient, err := DB.CreateIngredient(r.Context(), database.CreateIngredientParams{
		ID:        NewUUID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ParentID:  pgtype.UUID{Bytes: params.ParentID, Valid: params.ParentID != uuid.Nil},
	})
	if err != nil {
		log.Printf("error creating new ingredient: %s\n", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			respondError(w, http.StatusBadRequest, "ingredient already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondJSON(w, http.StatusOK, createIngredientResponse(ingredient))
}

func updateIngredientHandler(w http.ResponseWriter, r *http.Request) {
	ingredientIDString := chi.URLParam(r, "id")
	ingredientID, err := uuid.Parse(ingredientIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ingredient id not found")
		return
	}

	type parameters struct {
		Name     string    `json:"name"`
		ParentID uuid.UUID `json:"parent_id"`
	}
	params := &parameters{}
	err = json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	ingredient, err := DB.UpdateIngredientByID(r.Context(), database.UpdateIngredientByIDParams{
		ID:        pgUUID(ingredientID),
		Name:      params.Name,
		ParentID:  pgtype.UUID{Bytes: params.ParentID, Valid: params.ParentID != uuid.Nil},
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("error updating ingredient: %s\n", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			respondError(w, http.StatusBadRequest, "ingredient already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondJSON(w, http.StatusOK, createIngredientResponse(ingredient))
}

func listIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := DB.ListIngredients(r.Context())
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
