package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/views"
)

type RendererService interface {
	GetNavItems(isLoggedIn bool, currentURL string) []models.NavItem
}

func homeHandler(sm *scs.SessionManager, rds RendererService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := getUserIDFromCtx(r.Context(), sm)

		homeVM := views.NewHomeVM(userID, rds.GetNavItems(userID != uuid.Nil, r.URL.Path))
		render(w, r, views.Home(homeVM))
	}
}
