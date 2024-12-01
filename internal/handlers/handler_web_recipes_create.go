package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ajg/form"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-org/internal/models"
	views "github.com/quangd42/meal-org/internal/views/recipes"
)

func addRecipePageHandler(sm *scs.SessionManager, rds RendererService, rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCtx(r.Context(), sm)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodPost {
			rr, err := createMockRecipeRequest(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to mock recipe: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			recipe, err := rs.CreateRecipe(r.Context(), userID, rr)
			if err != nil {
				http.Error(w, "failed to create new recipe", http.StatusInternalServerError)
				return
			}

			render(w, r, views.RecipePostResponse(recipe.Name, true))
			return
		}

		vm := views.NewAddRecipeVM(userID, rds.GetNavItems(userID != uuid.Nil, r.URL.Path), nil)
		render(w, r, views.AddRecipePage(vm))
	}
}

func createMockRecipeRequest(r *http.Request) (models.RecipeRequest, error) {
	var rr models.RecipeRequest
	type request struct {
		Name        string `form:"name"`
		Description string `form:"description"`
		ExternalURL string `form:"external_url"`
	}
	arg := &request{}
	err := form.NewDecoder(r.Body).Decode(arg)
	if err != nil {
		return rr, err
	}
	if arg.Name == "" {
		return rr, errors.New("name is required")
	}
	rr = models.RecipeRequest{
		Name:              arg.Name,
		Description:       &arg.Description,
		ExternalURL:       &arg.ExternalURL,
		Servings:          0,
		CookTimeInMinutes: 0,
		Ingredients:       []models.IngredientInRecipe{},
		Instructions:      []models.InstructionInRecipe{},
	}
	return rr, err
}
