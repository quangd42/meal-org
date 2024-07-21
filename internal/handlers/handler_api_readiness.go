package handlers

import (
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	respondJSON(w, http.StatusOK, response{
		Status: http.StatusText(http.StatusOK),
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusInternalServerError, "Internal Server Error")
}
