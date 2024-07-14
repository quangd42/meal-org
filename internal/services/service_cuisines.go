package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/models"
)

func (rs RecipeService) CreateCuisine(ctx context.Context, cr models.CuisineRequest) (models.Cuisine, error) {
	var c models.Cuisine
	if cr.ParentID != nil {
		_, err := rs.store.Q.GetCuisineByID(ctx, *cr.ParentID)
		if err != nil {
			return c, checkErrNoRows(err)
		}
	}

	cuisine, err := rs.store.Q.CreateCuisine(ctx, database.CreateCuisineParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cr.Name,
		ParentID:  cr.ParentID,
	})
	if err != nil {
		return c, checkErrDBConstraint(err)
	}

	return createCuisineResponse(cuisine), nil
}

func (rs RecipeService) UpdateCuisineByID(ctx context.Context, cuisineID uuid.UUID, cr models.CuisineRequest) (models.Cuisine, error) {
	var c models.Cuisine
	if cr.ParentID != nil {
		_, err := rs.store.Q.GetCuisineByID(ctx, *cr.ParentID)
		if err != nil {
			return c, checkErrNoRows(err)
		}
	}

	cuisine, err := rs.store.Q.UpdateCuisineByID(ctx, database.UpdateCuisineByIDParams{
		ID:        cuisineID,
		Name:      cr.Name,
		ParentID:  cr.ParentID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return c, checkErrDBConstraint(err)
	}

	return createCuisineResponse(cuisine), nil
}

func (rs RecipeService) ListCuisines(ctx context.Context) ([]models.Cuisine, error) {
	var cs []models.Cuisine
	cuisines, err := rs.store.Q.ListCuisines(ctx)
	if err != nil {
		return cs, err
	}
	for _, dbc := range cuisines {
		cs = append(cs, createCuisineResponse(dbc))
	}
	return cs, nil
}

func (rs RecipeService) DeleteCuisine(ctx context.Context, cuisineID uuid.UUID) error {
	err := rs.store.Q.DeleteCuisine(ctx, cuisineID)
	if err != nil {
		return checkErrDBConstraint(err)
	}
	return nil
}

func createCuisineResponse(c database.Cuisine) models.Cuisine {
	res := models.Cuisine{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		ParentID:  c.ParentID,
	}
	return res
}
