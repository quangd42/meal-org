package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/quangd42/meal-planner/backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

var ErrAuthenticationFailed = errors.New("incorrect username or password")

func loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string
		Password string
	}
	params := &parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, ErrAuthenticationFailed.Error())
			return
		}

		respondInternalServerError(w)
		return
	}

	// compare password and hash
	err = auth.ValidateHash([]byte(user.Hash), []byte(params.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			respondError(w, http.StatusUnauthorized, ErrAuthenticationFailed.Error())
			return
		}
		respondInternalServerError(w)
		return
	}

	token, err := auth.CreateJWT(user, auth.DefaultExpirationDuration)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createResponseUserWithToken(user, token))
}
