package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

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

func respondError(w http.ResponseWriter, code int, value any) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s\n", value)
	}
	type response struct {
		Error any `json:"error"`
	}
	respondJSON(w, code, response{
		Error: value,
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

func respondMalformedRequestError(w http.ResponseWriter) {
	respondError(w, http.StatusBadRequest, "malformed request body")
}
