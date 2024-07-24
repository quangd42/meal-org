package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/views"
	"golang.org/x/crypto/bcrypt"
)

func loginPageHandler(sm *scs.SessionManager, rds RendererService, as AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			lr, err := decodeFormValidate[models.LoginRequest](r)
			if err != nil {
				vm := views.NewLoginVM(rds.GetNavItems(false, r.URL.Path), map[string]any{"email-password": true})
				render(w, r, views.LoginPage(vm))
				return
			}

			user, err := as.Login(r.Context(), lr)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					vm := views.NewLoginVM(rds.GetNavItems(false, r.URL.Path), map[string]any{"email-password": true})
					render(w, r, views.LoginPage(vm))
					return
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			sm.Put(r.Context(), "userID", user.ID)
			http.Redirect(w, r, fmt.Sprintf("http://%s/", r.Host), http.StatusSeeOther)
			return
		}
		vm := views.NewLoginVM(rds.GetNavItems(false, r.URL.Path), nil)
		render(w, r, views.LoginPage(vm))
	}
}
