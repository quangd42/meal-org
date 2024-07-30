package handlers

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/views"
)

func editRecipePageHandler(sm *scs.SessionManager, rds RendererService, rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCtx(r.Context(), sm)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		recipeIDStr := chi.URLParam(r, "recipeID")
		recipeID, err := uuid.Parse(recipeIDStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			rr, err := createMockRecipeRequest(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to mock recipe: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			recipe, err := rs.UpdateRecipeByID(r.Context(), userID, recipeID, rr)
			if err != nil {
				http.Error(w, "failed to update recipe", http.StatusInternalServerError)
				return
			}

			render(w, r, views.RecipePostResponse(recipe.Name))
			return
		}

		recipe, err := rs.GetRecipeByID(r.Context(), recipeID)
		if err != nil {
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		}

		vm := views.NewEditRecipeVM(userID, rds.GetNavItems(userID != uuid.Nil, r.URL.Path), recipe, nil)
		render(w, r, views.EditRecipePage(vm))
	}
}
