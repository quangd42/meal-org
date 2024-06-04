package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
	w.Write(data)
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
