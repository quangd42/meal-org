package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

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
