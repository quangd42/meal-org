package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
)

func AddRoutes(r *chi.Mux) {
	r.Get("/", r.NotFoundHandler())

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", readinessHandler)
		r.Get("/err", errorHandler)

		r.Mount("/users", usersAPIRouter())

		r.Post("/login", loginHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthVerifier())
			// r.Post("/logout", handleLogout)
		})
	})
}

// usersAPIRouter
func usersAPIRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", createUserHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthVerifier())

		r.Put("/", updateUserHandler)
	})

	return r
}
