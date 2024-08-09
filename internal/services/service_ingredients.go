package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/models"
)

func (rs RecipeService) CreateIngredient(ctx context.Context, arg models.IngredientRequest) (models.Ingredient, error) {
	var ing models.Ingredient

	ingredient, err := rs.store.Q.CreateIngredient(ctx, database.CreateIngredientParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      arg.Name,
	})
	if err != nil {
		return ing, checkErrDBConstraint(err)
	}

	ing = createIngredientResponse(ingredient)
	return ing, nil
}

func (rs RecipeService) UpdateIngredientByID(ctx context.Context, ingredientID uuid.UUID, arg models.IngredientRequest) (models.Ingredient, error) {
	var ing models.Ingredient

	ingredient, err := rs.store.Q.UpdateIngredientByID(ctx, database.UpdateIngredientByIDParams{
		ID:        ingredientID,
		Name:      arg.Name,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return ing, customDBErr(err)
	}

	ing = createIngredientResponse(ingredient)
	return ing, nil
}

func (rs RecipeService) ListIngredients(ctx context.Context) ([]models.Ingredient, error) {
	var ings []models.Ingredient
	ingredients, err := rs.store.Q.ListIngredients(ctx)
	if err != nil {
		return ings, err
	}

	for _, ing := range ingredients {
		ings = append(ings, createIngredientResponse(ing))
	}
	return ings, err
}

func (rs RecipeService) DeleteIngredient(ctx context.Context, ingredientID uuid.UUID) error {
	err := rs.store.Q.DeleteIngredient(ctx, ingredientID)
	if err != nil {
		return checkErrDBConstraint(err)
	}
	return nil
}

func createIngredientResponse(i database.Ingredient) models.Ingredient {
	res := models.Ingredient{
		ID:        i.ID,
		Name:      i.Name,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
	return res
}
