package handlers

import (
	"github.com/go-chi/chi/v5"
)

func AddRoutes(r *chi.Mux) {
	r.Get("/", r.NotFoundHandler())

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", HandleReadiness)
		r.Get("/err", HandleError)
	})
}
