package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/quangd42/meal-org/internal/auth"
	"github.com/quangd42/meal-org/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrAuthenticationFailed = errors.New("incorrect email or password")

type AuthService interface {
	GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, error)
	GenerateAndSaveRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	Login(ctx context.Context, lr models.LoginRequest) (models.User, error)
	AuthVerifier() func(http.Handler) http.Handler
}

func loginAPIHandler(as AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lr, err := decodeJSONValidate[models.LoginRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		user, err := as.Login(r.Context(), lr)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				respondError(w, http.StatusUnauthorized, ErrAuthenticationFailed.Error())
				return
			}

			respondInternalServerError(w)
			return
		}

		jwt, err := as.GenerateAccessToken(r.Context(), user.ID)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		refreshToken, err := as.GenerateAndSaveRefreshToken(r.Context(), user.ID)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, user.WithToken(jwt, refreshToken))
	}
}

func refreshAccessHandler(as AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetHeaderToken(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, auth.ErrTokenNotFound.Error())
			return
		}

		userID, err := as.ValidateRefreshToken(r.Context(), refreshToken)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		jwt, err := as.GenerateAccessToken(r.Context(), userID)
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
}

func revokeRefreshTokenHandler(as AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetHeaderToken(r)
		if err != nil {
			http.Error(w, auth.ErrTokenNotFound.Error(), http.StatusUnauthorized)
			return
		}

		err = as.RevokeRefreshToken(r.Context(), refreshToken)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
