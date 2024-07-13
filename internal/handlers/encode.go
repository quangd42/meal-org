package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

func respondJSON[T any](w http.ResponseWriter, code int, v T) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("error decoding JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data) // #nosec G104
}

func respondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s\n", msg)
	}
	type response struct {
		Error string `json:"error"`
	}
	respondJSON(w, code, response{
		Error: msg,
	})
}

func respondInternalServerError(w http.ResponseWriter) {
	respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func respondDBConstraintsError(w http.ResponseWriter, err error, msg string) {
	if errors.Is(err, services.ErrDBConstraint) {
		respondError(w, http.StatusForbidden, fmt.Sprintf("invalid operation, check: %s", msg))
		return
	}
	respondInternalServerError(w)
}

func respondUniqueValueError(w http.ResponseWriter, err error, msg string) {
	if errors.Is(err, services.ErrUniqueValue) {
		respondError(w, http.StatusBadRequest, "unique value constraint violated: "+msg)
		return
	}
	respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func respondMalformedRequestError(w http.ResponseWriter) {
	respondError(w, http.StatusBadRequest, "malformed request body")
}

func createIngredientResponse(i database.Ingredient) models.Ingredient {
	res := models.Ingredient{
		ID:        i.ID,
		Name:      i.Name,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		ParentID:  i.ParentID,
	}
	return res
}

func createCuisineResponse(i database.Cuisine) models.Cuisine {
	res := models.Cuisine{
		ID:        i.ID,
		Name:      i.Name,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		ParentID:  i.ParentID,
	}
	return res
}
