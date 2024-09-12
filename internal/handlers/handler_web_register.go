package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/quangd42/meal-planner/internal/models"
	"github.com/quangd42/meal-planner/internal/models/validator"
	"github.com/quangd42/meal-planner/internal/services"
	views "github.com/quangd42/meal-planner/internal/views/auth"
)

func registerPageHandler(sm *scs.SessionManager, rds RendererService, us UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			errs := make(map[string]any)
			ur, err := decodeFormValidate[models.CreateUserRequest](r)
			if err != nil {
				for errName, errMsg := range err.(validator.ValidationErrors) {
					// errName := strings.ToLower(fmt.Sprintf("%s-%s", err.Field(), err.Tag()))
					fmt.Println(errName)
					errs[errName] = errMsg
				}
				render(w, r, views.RegisterForm(errs))
				return
			}

			user, err := us.CreateUser(r.Context(), ur)
			if err != nil {
				if errors.Is(err, services.ErrHashPassword) {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				errs["email-duplicate"] = true
				render(w, r, views.RegisterForm(errs))
				return
			}

			sm.Put(r.Context(), "userID", user.ID)
			w.Header().Set("HX-Redirect", fmt.Sprintf("http://%s/", r.Host))
			return
		}
		vm := views.NewRegisterVM(rds.GetNavItems(false, r.URL.Path), nil)
		render(w, r, views.RegisterPage(vm))
	}
}
