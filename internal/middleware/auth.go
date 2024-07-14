package middleware

import (
	"context"
	"net/http"

	"github.com/quangd42/meal-planner/internal/services/auth"
)

type contextKey struct {
	name string
}

var (
	UserIDCtxKey = &contextKey{"userID"}
	TokenCtxKey  = &contextKey{"token"}
)

// TODO: split this into two: on to verify if token is good, one to
// extract, verify and return user information
func AuthVerifier() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.GetHeaderToken(r)
			if err != nil {
				http.Error(w, auth.ErrTokenNotFound.Error(), http.StatusUnauthorized)
				return
			}
			userID, err := auth.VerifyJWT(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDCtxKey, userID)
			ctx = context.WithValue(ctx, TokenCtxKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
