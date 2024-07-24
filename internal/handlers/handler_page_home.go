package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/views"
)

type RendererService interface {
	GetNavItems(isLoggedIn bool, currentURL string) []views.NavItem
}

func homeHandler(sm *scs.SessionManager, rds RendererService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCtx(r.Context(), sm)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		homeVM := views.NewHomeVM(userID, rds.GetNavItems(userID != uuid.Nil, r.URL.Path))
		render(w, r, views.Home(homeVM))
	}
}
