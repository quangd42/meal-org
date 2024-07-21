package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/components"
	"github.com/quangd42/meal-planner/internal/views"
)

var navItems = []components.NavItem{
	{
		Name: "Recipes",
		URL:  templ.URL("/"),
	},
	{
		Name: "Add a Recipe",
		URL:  templ.URL("/recipes/add"),
	},
}

func homeHandler(sm *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := sm.Get(r.Context(), "userID").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		homeVM := views.NewHomeVM(userID, navItems)
		views.Home(homeVM).Render(r.Context(), w)
	}
}
