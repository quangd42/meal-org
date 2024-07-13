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

func createCuisineHandler(w http.ResponseWriter, r *http.Request) {
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
		_, err = store.Q.GetCuisineByID(r.Context(), *params.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent cuisine does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
	}

	cuisine, err := store.Q.CreateCuisine(r.Context(), database.CreateCuisineParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ParentID:  params.ParentID,
	})
	if err != nil {
		log.Printf("error creating new cuisine: %s\n", err)
		respondUniqueValueError(w, err, "cuisine name")
		return
	}

	respondJSON(w, http.StatusCreated, createCuisineResponse(cuisine))
}

func updateCuisineHandler(w http.ResponseWriter, r *http.Request) {
	cuisineIDString := chi.URLParam(r, "id")
	cuisineID, err := uuid.Parse(cuisineIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "cuisine id not found")
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
		_, err = store.Q.GetCuisineByID(r.Context(), *params.ParentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respondError(w, http.StatusBadRequest, "parent cuisine does not exist")
				return
			}
			respondInternalServerError(w)
			return
		}
	}

	cuisine, err := store.Q.UpdateCuisineByID(r.Context(), database.UpdateCuisineByIDParams{
		ID:        cuisineID,
		Name:      params.Name,
		ParentID:  params.ParentID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("error updating cuisine: %s\n", err)
		respondUniqueValueError(w, err, "cuisine name")
		return
	}

	respondJSON(w, http.StatusOK, createCuisineResponse(cuisine))
}

func listCuisinesHandler(w http.ResponseWriter, r *http.Request) {
	cuisines, err := store.Q.ListCuisines(r.Context())
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "no cuisines found")
			return
		}
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, cuisines)
}

func deleteCuisineHandler(w http.ResponseWriter, r *http.Request) {
	cuisineIDString := chi.URLParam(r, "id")
	cuisineID, err := uuid.Parse(cuisineIDString)
	if err != nil {
		respondError(w, http.StatusBadRequest, "cuisine id not found")
		return
	}

	err = store.Q.DeleteCuisine(r.Context(), cuisineID)
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
