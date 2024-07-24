package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/quangd42/meal-planner/internal/middleware"
)

func AddRoutes(
	r *chi.Mux,
	sm *scs.SessionManager,
	rds RendererService,
	us UserService,
	as AuthService,
	rs RecipeService,
	is IngredientService,
	cs CuisineService,
) {
	// Top level middlewares
	r.Use(chiMiddleware.StripSlashes)
	r.Use(sm.LoadAndSave)

	// Static assets
	fs := disableCacheInDevMode(http.FileServer(http.Dir("assets")))
	r.Handle("/assets/*", http.StripPrefix("/assets", fs))

	// Public pages
	r.Get("/login", loginPageHandler(sm, rds, as))
	r.Post("/login", loginPageHandler(sm, rds, as))
	r.Post("/logout", logoutHandler(sm))
	r.Get("/register", registerPageHandler(sm, rds, us))
	r.Post("/register", registerPageHandler(sm, rds, us))
	r.Get("/", homeHandler(sm, rds))

	// Private pages
	r.Get("/recipes/add", addRecipePageHandler(sm, rds))

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", readinessHandler)
		r.Get("/err", errorHandler)

		r.Mount("/users", usersAPIRouter(us, as))
		r.Mount("/auth", authAPIRouter(as))
		r.Mount("/recipes", recipesAPIRouter(rs))
		r.Mount("/ingredients", ingredientsAPIRouter(is))
		r.Mount("/cuisines", cuisinesAPIRouter(cs))
	})
}

// authRouter
func authRouter(as AuthService) http.Handler {
	r := chi.NewRouter()

	_ = as

	return r
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

// authAPIRouter
func authAPIRouter(as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", loginAPIHandler(as))
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
func cuisinesAPIRouter(cs CuisineService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AuthVerifier())
	r.Post("/", createCuisineHandler(cs))
	r.Get("/", listCuisinesHandler(cs))

	r.Put("/{id}", updateCuisineHandler(cs))
	r.Delete("/{id}", deleteCuisineHandler(cs))

	return r
}
