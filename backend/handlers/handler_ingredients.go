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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

// TODO: each user can create their own ingredients
// admin create shared set
func createIngredientHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string    `json:"name"`
		ParentID uuid.UUID `json:"parent_id"`
	}
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		respondMalformedRequestError(w)
		return
	}

	var parentID pgtype.UUID
	if params.ParentID == uuid.Nil {
		parentID = pgtype.UUID{Valid: false}
	} else {
		_, err = DB.GetIngredientByID(r.Context(), pgUUID(params.ParentID))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent ingredient does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
		parentID = pgUUID(params.ParentID)
	}

	ingredient, err := DB.CreateIngredient(r.Context(), database.CreateIngredientParams{
		ID:        NewUUID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ParentID:  parentID,
	})
	if err != nil {
		log.Printf("error creating new ingredient: %s\n", err)
		respondUniqueValueError(w, err, "ingredient name")
		return
	}

	respondJSON(w, http.StatusCreated, createIngredientResponse(ingredient))
}

// TODO: each user can update their own ingredients
// admin can update the shared set
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
		respondMalformedRequestError(w)
		return
	}

	var parentID pgtype.UUID
	if params.ParentID == uuid.Nil {
		parentID = pgtype.UUID{Valid: false}
	} else {
		_, err = DB.GetIngredientByID(r.Context(), pgUUID(params.ParentID))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent ingredient does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
		parentID = pgUUID(params.ParentID)
	}

	ingredient, err := DB.UpdateIngredientByID(r.Context(), database.UpdateIngredientByIDParams{
		ID:        pgUUID(ingredientID),
		Name:      params.Name,
		ParentID:  parentID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("error updating ingredient: %s\n", err)
		respondUniqueValueError(w, err, "ingredient name")
		return
	}

	respondJSON(w, http.StatusOK, createIngredientResponse(ingredient))
}

// TODO: each user should see their own set
func listIngredientsHandler(w http.ResponseWriter, r *http.Request) {
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

// TODO: normal user can only delete their own ingredients
func deleteIngredientHandler(w http.ResponseWriter, r *http.Request) {
	ingredientIDString := chi.URLParam(r, "id")
	ingredientID, err := uuid.Parse(ingredientIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ingredient id not found")
		return
	}

	err = DB.DeleteIngredient(r.Context(), pgUUID(ingredientID))
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
