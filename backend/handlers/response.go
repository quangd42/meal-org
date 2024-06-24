package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/backend/internal/database"
)

func respondJSON[T any](w http.ResponseWriter, code int, v T) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("error decoding JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s\n", msg)
	}
	type response struct {
		Error string `json:"error"`
	}
	respondJSON(w, code, response{
		Error: msg,
	})
}

func respondInternalServerError(w http.ResponseWriter) {
	respondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Hash         string    `json:"-"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func createUserResponse(u database.User) User {
	return User{
		ID:        u.ID.Bytes,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Name:      u.Name,
		Username:  u.Username,
	}
}

func createUserResponseWithToken(u database.User, token, refreshToken string) User {
	user := createUserResponse(u)
	user.Token = token
	user.RefreshToken = refreshToken
	return user
}

type Recipe struct {
	ID          uuid.UUID            `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Name        string               `json:"name"`
	ExternalUrl string               `json:"external_url"`
	UserID      uuid.UUID            `json:"user_id"`
	Ingredients []IngredientInRecipe `json:"ingredients"`
}

type IngredientInRecipe struct {
	ID          uuid.UUID  `json:"id"`
	Amount      string     `json:"amount"`
	Instruction string     `json:"instruction"`
	Name        string     `json:"name"`
	ParentID    *uuid.UUID `json:"parent_id"`
}

func createRecipeResponse(dr database.Recipe, dis []database.ListIngredientsByRecipeIDRow) Recipe {
	ingredients := []IngredientInRecipe{}
	for _, di := range dis {
		ingredients = append(ingredients, IngredientInRecipe{
			ID:          di.ID.Bytes,
			Amount:      di.Amount,
			Instruction: di.Instruction.String,
			Name:        di.Name,
		})
	}
	return Recipe{
		ID:          dr.ID.Bytes,
		CreatedAt:   dr.CreatedAt,
		UpdatedAt:   dr.UpdatedAt,
		Name:        dr.Name,
		ExternalUrl: dr.ExternalUrl,
		UserID:      dr.UserID.Bytes,
		Ingredients: ingredients,
	}
}

type Ingredient struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"`
}

func createIngredientResponse(i database.Ingredient) Ingredient {
	res := Ingredient{
		ID:        i.ID.Bytes,
		Name:      i.Name,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		ParentID:  &uuid.UUID{},
	}
	if i.ParentID.Valid {
		*res.ParentID = i.ParentID.Bytes
	} else {
		res.ParentID = nil
	}
	return res
}
