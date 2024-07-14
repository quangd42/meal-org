package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/quangd42/meal-planner/internal/middleware"
)

func AddRoutes(r *chi.Mux,
	us UserService,
	as AuthService,
	rs RecipeService,
	is IngredientService,
) {
	r.Get("/", r.NotFoundHandler())

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", readinessHandler)
		r.Get("/err", errorHandler)

		r.Mount("/users", usersAPIRouter(us, as))
		r.Mount("/auth", authRouter(as))
		r.Mount("/recipes", recipesAPIRouter(rs))
		r.Mount("/ingredients", ingredientsAPIRouter(is))
		r.Mount("/cuisines", cuisinesAPIRouter())
	})
}

// usersAPIRouter
func usersAPIRouter(us UserService, as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Post("/", createUserHandler(us, as))
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthVerifier())
		r.Put("/", updateUserHandler(us))
		r.Delete("/", forgetMeHandler(us))
	})

	return r
}

// authRouter
func authRouter(as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", loginHandler(as))
	r.Post("/refresh", refreshAccessHandler(as))
	r.Post("/revoke", revokeRefreshTokenHandler(as))

	return r
}

// recipesAPIRouter
func recipesAPIRouter(rs RecipeService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AuthVerifier())
	r.Post("/", createRecipeHandler(rs))
	r.Get("/", listRecipesHandler(rs))

	r.Get("/{id}", getRecipeHandler(rs))
	r.Put("/{id}", updateRecipeHandler(rs))
	r.Delete("/{id}", deleteRecipeHandler(rs))

	// TODO: add search & filter

	return r
}

func ingredientsAPIRouter(is IngredientService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AuthVerifier())
	r.Post("/", createIngredientHandler(is))
	r.Get("/", listIngredientsHandler(is))

	r.Put("/{id}", updateIngredientHandler(is))
	r.Delete("/{id}", deleteIngredientHandler(is))

	return r
}

// TODO: create, update and delete should be restricted to admin only
func cuisinesAPIRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AuthVerifier())
	r.Post("/", createCuisineHandler)
	r.Get("/", listCuisinesHandler)

	r.Put("/{id}", updateCuisineHandler)
	r.Delete("/{id}", deleteCuisineHandler)

	return r
}
