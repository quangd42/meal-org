package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func deleteRecipePageHandler(sm *scs.SessionManager, rs RecipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := getUserIDFromCtx(r.Context(), sm)
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

		err = rs.DeleteRecipeByID(r.Context(), recipeID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
