package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	views "github.com/quangd42/meal-org/internal/views/recipes"
)

func listRecipesPageHandler(sm *scs.SessionManager, rds RendererService, rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCtx(r.Context(), sm)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		recipes, err := rs.ListRecipesWithCuisinesByUserID(r.Context(), userID, getPaginationParams(r))
		if err != nil {
			http.Error(w, "internal error", 500)
			return
		}
		render(w, r, views.ListRecipesPage(views.NewListRecipesVM(rds.GetNavItems(true, r.URL.Path), recipes, nil)))
	}
}
