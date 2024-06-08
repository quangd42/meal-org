package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AddRoutes(r *chi.Mux, c *Config) {
	r.Get("/", r.NotFoundHandler())

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", HandleReadiness)
		r.Get("/err", HandleError)

		r.Mount("/users", usersAPIRouter(c))
	})
}

// usersAPIRouter
func usersAPIRouter(c *Config) http.Handler {
	r := chi.NewRouter()

	r.Post("/", CreateUserHandler(c))

	return r
}
