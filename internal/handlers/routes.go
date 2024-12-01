package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func AddRoutes(
	r *chi.Mux,
	sm *scs.SessionManager,
	rds RendererService,
	us UserService,
	as AuthService,
	rs RecipeService,
) {
	// Top level middlewares
	r.Use(middleware.StripSlashes)
	r.Use(sm.LoadAndSave)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	// Add
	r.Get("/recipes/add", addRecipePageHandler(sm, rds, rs))
	r.Post("/recipes", addRecipePageHandler(sm, rds, rs))
	// List
	r.Get("/recipes", listRecipesPageHandler(sm, rds, rs))
	// Edit
	r.Post("/recipes/{recipeID}", editRecipePageHandler(sm, rds, rs))
	r.Get("/recipes/{recipeID}", editRecipePageHandler(sm, rds, rs))
	// Delete
	r.Delete("/recipes/{recipeID}", deleteRecipePageHandler(sm, rs))

	// API router
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthz", readinessHandler)
		r.Get("/err", errorHandler)

		r.Mount("/users", usersAPIRouter(us, as))
		r.Mount("/auth", authAPIRouter(as))
		r.Mount("/recipes", recipesAPIRouter(rs, as))
		r.Mount("/ingredients", ingredientsAPIRouter(rs, as))
		r.Mount("/cuisines", cuisinesAPIRouter(rs, as))
	})
}

// usersAPIRouter
func usersAPIRouter(us UserService, as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Post("/", createUserHandler(us, as))
	r.Group(func(r chi.Router) {
		r.Use(as.AuthVerifier())
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
func recipesAPIRouter(rs RecipeService, as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(as.AuthVerifier())
	r.Post("/", createRecipeHandler(rs))
	r.Get("/", listRecipesHandler(rs))

	r.Get("/{id}", getRecipeHandler(rs))
	r.Put("/{id}", updateRecipeHandler(rs))
	r.Delete("/{id}", deleteRecipeHandler(rs))

	// TODO: add search & filter

	return r
}

func ingredientsAPIRouter(rs RecipeService, as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(as.AuthVerifier())
	r.Post("/", createIngredientHandler(rs))
	r.Get("/", listIngredientsHandler(rs))

	r.Put("/{id}", updateIngredientHandler(rs))
	r.Delete("/{id}", deleteIngredientHandler(rs))

	return r
}

// TODO: create, update and delete should be restricted to admin only
func cuisinesAPIRouter(rs RecipeService, as AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(as.AuthVerifier())
	r.Post("/", createCuisineHandler(rs))
	r.Get("/", listCuisinesHandler(rs))

	r.Put("/{id}", updateCuisineHandler(rs))
	r.Delete("/{id}", deleteCuisineHandler(rs))

	return r
}
