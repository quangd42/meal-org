package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/quangd42/meal-planner/backend/internal/database"
	"github.com/quangd42/meal-planner/backend/internal/models"
)

type CreateWholeRecipeParams struct {
	UserID uuid.UUID `json:"user_id"`
	models.RecipeRequest
}

func CreateWholeRecipe(ctx context.Context, store *database.Store, arg CreateWholeRecipeParams) (models.Recipe, error) {
	var r models.Recipe

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return models.Recipe{}, err
	}
	defer tx.Rollback(ctx)

	qtx := store.Q.WithTx(tx)

	// Create host Recipe
	dbRecipe, err := qtx.CreateRecipe(ctx, database.CreateRecipeParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		Name:              arg.Name,
		ExternalUrl:       arg.ExternalURL,
		Servings:          int32(arg.Servings),
		Yield:             arg.Yield,
		CookTimeInMinutes: int32(arg.CookTimeInMinutes),
		Notes:             arg.Notes,
		UserID:            arg.UserID,
	})
	if err != nil {
		return r, err
	}

	// Add Ingredients to host Recipe
	dbIngredientParams := make([]database.AddIngredientsToRecipeParams, len(arg.Ingredients))
	for i, p := range arg.Ingredients {
		dbIngredientParams[i] = database.AddIngredientsToRecipeParams{
			RecipeID:     dbRecipe.ID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			IngredientID: p.ID,
			Amount:       p.Amount,
			PrepNote:     p.PrepNote,
		}
	}

	// Can add as many ingredients as desired, no check here
	_, err = qtx.AddIngredientsToRecipe(ctx, dbIngredientParams)
	if err != nil {
		return r, err
	}

	dbIngredients, err := qtx.ListIngredientsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	// Add Instructions to host Recipe
	for _, p := range arg.Instructions {
		err = qtx.AddInstructionToRecipe(ctx, database.AddInstructionToRecipeParams{
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			StepNo:      int32(p.StepNo),
			Instruction: p.Instruction,
			RecipeID:    dbRecipe.ID,
		})
		if err != nil {
			return r, err
		}
	}

	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	r = createWholeRecipe(dbRecipe, dbIngredients, dbInstructions)

	return r, tx.Commit(ctx)
}

type UpdateWholeRecipeParams struct {
	ID uuid.UUID `json:"id"`
	models.RecipeRequest
}

func UpdateWholeRecipe(ctx context.Context, store *database.Store, arg UpdateWholeRecipeParams) (models.Recipe, error) {
	var r models.Recipe
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return r, err
	}
	defer tx.Rollback(ctx)

	qtx := store.Q.WithTx(tx)

	// Update host Recipe
	dbRecipe, err := qtx.UpdateRecipeByID(ctx, database.UpdateRecipeByIDParams{
		ID:                arg.ID,
		UpdatedAt:         time.Now().UTC(),
		Name:              arg.Name,
		ExternalUrl:       arg.ExternalURL,
		Servings:          int32(arg.Servings),
		Yield:             arg.Yield,
		CookTimeInMinutes: int32(arg.CookTimeInMinutes),
		Notes:             arg.Notes,
	})
	if err != nil {
		return r, err
	}

	// Update ingredients in Recipe

	// Make sure that ingredients already belong to the host recipe
	dbIngredients, err := qtx.ListIngredientsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}
	dbIngredientsMap := make(map[uuid.UUID]bool, len(dbIngredients))
	for _, i := range dbIngredients {
		dbIngredientsMap[i.ID] = true
	}

	for _, i := range arg.Ingredients {
		if _, ok := dbIngredientsMap[i.ID]; !ok {
			pgErr := &pgconn.PgError{
				Code: "23000",
			}
			return r, pgErr
		}
		ingreDBParams := database.UpdateIngredientInRecipeParams{
			Amount:       i.Amount,
			PrepNote:     i.PrepNote,
			UpdatedAt:    time.Now().UTC(),
			IngredientID: i.ID,
			RecipeID:     dbRecipe.ID,
		}

		err = qtx.UpdateIngredientInRecipe(ctx, ingreDBParams)
		if err != nil {
			return r, err
		}
	}

	dbIngredients, err = qtx.ListIngredientsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	// Update instructions in Recipe
	// List instructions from db
	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	// Make map from db
	dbInstructionsMap := make(map[int]database.Instruction, len(dbInstructions))
	for _, dbi := range dbInstructions {
		dbInstructionsMap[int(dbi.StepNo)] = dbi
	}

	var toAdd, toUpdate []models.InstructionInRecipe
	for _, pi := range arg.Instructions {
		_, ok := dbInstructionsMap[pi.StepNo]
		// If in param the step no is found in db, add it to the update list
		// then delete it from the map
		if ok {
			toUpdate = append(toUpdate, pi)
			delete(dbInstructionsMap, pi.StepNo)
			// Else the param is new, add it to the add list
		} else {
			toAdd = append(toAdd, pi)
		}
	}

	for _, pi := range toAdd {
		err = qtx.AddInstructionToRecipe(ctx, database.AddInstructionToRecipeParams{
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Instruction: pi.Instruction,
			StepNo:      int32(pi.StepNo),
			RecipeID:    dbRecipe.ID,
		})
		if err != nil {
			return r, err
		}
	}

	for _, pi := range toUpdate {
		err = qtx.UpdateInstructionByID(ctx, database.UpdateInstructionByIDParams{
			UpdatedAt:   time.Now().UTC(),
			Instruction: pi.Instruction,
			StepNo:      int32(pi.StepNo),
			RecipeID:    dbRecipe.ID,
		})
		if err != nil {
			return r, err
		}
	}

	// The rest of the dbMap is not found in param, so delete them
	for _, dbi := range dbInstructionsMap {
		err = qtx.DeleteInstructionByID(ctx, database.DeleteInstructionByIDParams{
			StepNo:   dbi.StepNo,
			RecipeID: dbi.RecipeID,
		})
		if err != nil {
			return r, err
		}
	}

	dbInstructions, err = qtx.ListInstructionsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	// Assemble all updated data
	r = createWholeRecipe(dbRecipe, dbIngredients, dbInstructions)

	return r, tx.Commit(ctx)
}

func GetWholeRecipe(ctx context.Context, store *database.Store, recipeID uuid.UUID) (models.Recipe, error) {
	var r models.Recipe

	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return models.Recipe{}, err
	}
	defer tx.Rollback(ctx)

	qtx := store.Q.WithTx(tx)

	dbRecipe, err := qtx.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return r, err
	}

	ingredients, err := qtx.ListIngredientsByRecipeID(ctx, recipeID)
	if err != nil {
		return r, err
	}

	instructions, err := qtx.ListInstructionsByRecipeID(ctx, recipeID)
	if err != nil {
		return r, err
	}
	r = createWholeRecipe(dbRecipe, ingredients, instructions)
	return r, tx.Commit(ctx)
}

func createWholeRecipe(dr database.Recipe, dbIngredients []database.ListIngredientsByRecipeIDRow, dbInstructions []database.Instruction) models.Recipe {
	ingredients := []models.IngredientInRecipe{}
	for _, di := range dbIngredients {
		ingredients = append(ingredients, models.IngredientInRecipe{
			ID:       di.ID,
			Amount:   di.Amount,
			PrepNote: di.PrepNote,
			Name:     di.Name,
		})
	}

	instructions := []models.InstructionInRecipe{}
	for _, di := range dbInstructions {
		instructions = append(instructions, models.InstructionInRecipe{
			StepNo:      int(di.StepNo),
			Instruction: di.Instruction,
		})
	}

	return models.Recipe{
		ID:                dr.ID,
		CreatedAt:         dr.CreatedAt,
		UpdatedAt:         dr.UpdatedAt,
		Name:              dr.Name,
		ExternalUrl:       dr.ExternalUrl,
		UserID:            dr.UserID,
		Servings:          int(dr.Servings),
		Yield:             dr.Yield,
		CookTimeInMinutes: int(dr.CookTimeInMinutes),
		Notes:             dr.Notes,
		Ingredients:       ingredients,
		Instructions:      instructions,
	}
}
