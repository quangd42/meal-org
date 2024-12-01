package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/quangd42/meal-org/internal/auth"
)

type contextKey int

const (
	_ contextKey = iota
	userIDCtxKey
	tokenCtxKey
)

// TODO: split this into two: on to verify if token is good, one to
// extract, verify and return user information

func (as Auth) AuthVerifier() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetHeaderToken(r)
			if err != nil {
				http.Error(w, auth.ErrTokenNotFound.Error(), http.StatusUnauthorized)
				return
			}
			userID, err := auth.VerifyJWT(as.jwtSecret, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = ContextWithUserID(ctx, userID)
			ctx = ContextWithToken(ctx, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func UserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value(userIDCtxKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, nil
	}
	return userID, nil
}

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDCtxKey, userID)
}

func TokenFromContext(r *http.Request) (string, error) {
	token, ok := r.Context().Value(tokenCtxKey).(string)
	if !ok {
		return "", nil
	}
	return token, nil
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}
