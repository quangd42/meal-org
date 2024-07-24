package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/views"
)

func addRecipePageHandler(sm *scs.SessionManager, rds RendererService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := sm.Get(r.Context(), "userID").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		homeVM := views.NewHomeVM(userID, rds.GetNavItems(userID != uuid.Nil))
		views.Home(homeVM).Render(r.Context(), w)
	}
}
