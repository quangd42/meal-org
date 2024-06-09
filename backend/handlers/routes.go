package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/quangd42/meal-planner/backend/internal/middleware"
)

func AddRoutes(r *chi.Mux, c *Config) {
	r.Get("/", r.NotFoundHandler())

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", HandleReadiness)
		r.Get("/err", HandleError)

		r.Mount("/users", usersAPIRouter(c))

		r.Group(func(r chi.Router) {
			r.Post("/login", handleLogin)
		})
	})
}

// usersAPIRouter
func usersAPIRouter(c *Config) http.Handler {
	r := chi.NewRouter()

	r.Post("/", CreateUserHandler(c))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthVerifier())

		r.Put("/", UpdateUserHandler(c))
	})

	return r
}
