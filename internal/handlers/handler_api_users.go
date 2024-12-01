package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/auth"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/services"
)

type UserService interface {
	CreateUser(ctx context.Context, ur models.CreateUserRequest) (models.User, error)
	// GetUserByID() (models.User, error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, ur models.UpdateUserRequest) (models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
}

func createUserHandler(us UserService, as AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ur, err := decodeJSONValidate[models.CreateUserRequest](r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		user, err := us.CreateUser(r.Context(), ur)
		if err != nil {
			if errors.Is(err, services.ErrHashPassword) {
				respondInternalServerError(w)
				return
			}
			if errors.Is(err, services.ErrDBConstraint) {
				respondError(w, http.StatusForbidden, "email already taken")
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

		respondJSON(w, http.StatusCreated, user.WithToken(jwt, refreshToken))
	}
}

func updateUserHandler(us UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.UserIDFromContext(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
			return
		}

		ur, err := decodeJSONValidate[models.UpdateUserRequest](r)
		if err != nil {
			respondMalformedRequestError(w)
			return
		}

		user, err := us.UpdateUserByID(r.Context(), userID, ur)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, user)
	}
}

func forgetMeHandler(us UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.UserIDFromContext(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, auth.ErrTokenNotFound.Error())
			return
		}

		err = us.DeleteUserByID(r.Context(), userID)
		if err != nil {
			respondInternalServerError(w)
			return
		}

		respondJSON(w, http.StatusNoContent, "user deleted")
	}
}
