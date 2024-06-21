package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

func createIngredientHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string    `json:"name"`
		ParentID uuid.UUID `json:"parent_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	ingredient, err := DB.CreateIngredient(r.Context(), database.CreateIngredientParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ParentID:  params.ParentID,
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

	respondJSON(w, http.StatusOK, ingredient)
}
