package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/quangd42/meal-planner/backend/internal/auth"
	"github.com/quangd42/meal-planner/backend/internal/database"
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
		if errors.Is(err, pgx.ErrNoRows) {
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

	jwt, refreshToken, err := generateAndSaveAuthTokens(r, user)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	respondJSON(w, http.StatusOK, createUserResponseWithToken(user, jwt, refreshToken))
}

// TODO: requires authentication for refreshing and revoking tokens
func refreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	paramToken, err := auth.GetHeaderToken(r)
	if err != nil {
		http.Error(w, auth.ErrTokenNotFound.Error(), http.StatusUnauthorized)
		return
	}

	token, err := DB.GetTokenByValue(r.Context(), paramToken)
	if err != nil {
		respondError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if token.ExpiredAt.Before(time.Now().UTC()) {
		token.IsRevoked = true
		err = DB.RevokeToken(r.Context(), database.RevokeTokenParams{
			Value:     paramToken,
			IsRevoked: true,
		})
		if err != nil {
			respondInternalServerError(w)
			return
		}
	}
	if token.IsRevoked {
		respondError(w, http.StatusUnauthorized, auth.ErrTokenInvalid.Error())
		return
	}

	jwt, err := auth.CreateJWT(database.User{ID: token.UserID}, auth.ExpirationDurationAccess)
	if err != nil {
		respondInternalServerError(w)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}

func revokeJWTHandler(w http.ResponseWriter, r *http.Request) {
	paramToken, err := auth.GetHeaderToken(r)
	if err != nil {
		http.Error(w, auth.ErrTokenNotFound.Error(), http.StatusUnauthorized)
		return
	}

	_, err = DB.GetTokenByValue(r.Context(), paramToken)
	if err != nil {
		respondError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	err = DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Value:     paramToken,
		IsRevoked: true,
	})
	if err != nil {
		respondInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
