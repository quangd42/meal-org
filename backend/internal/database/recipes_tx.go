package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/backend/internal/models"
)

type step struct {
	StepNo      int    `json:"step_no"`
	Instruction string `json:"instruction"`
}

type UpdateAllRecipeDataParams struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	ExternalURL       *string   `json:"external_url"`
	Servings          int       `json:"servings"`
	Yield             *string   `json:"yield"`
	CookTimeInMinutes int       `json:"cook_time_in_minutes"`
	Notes             *string   `json:"notes"`
	Ingredients       []struct {
		ID       uuid.UUID `json:"id"`
		Amount   string    `json:"amount"`
		PrepNote *string   `json:"prep_note"`
	} `json:"ingredients"`
	Instructions []step `json:"instructions"`
}

func UpdateAllRecipeData(ctx context.Context, store *Store, arg UpdateAllRecipeDataParams) (models.Recipe, error) {
	var r models.Recipe
	tx, err := store.DB.Begin(ctx)
	if err != nil {
		return r, err
	}
	defer tx.Rollback(ctx)

	qtx := store.Q.WithTx(tx)

	// Update host Recipe
	recipe, err := qtx.UpdateRecipeByID(ctx, UpdateRecipeByIDParams{
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
	for _, i := range arg.Ingredients {
		ingreDBParams := UpdateIngredientInRecipeParams{
			Amount:       i.Amount,
			PrepNote:     i.PrepNote,
			UpdatedAt:    time.Now().UTC(),
			IngredientID: i.ID,
			RecipeID:     arg.ID,
		}

		err = qtx.UpdateIngredientInRecipe(ctx, ingreDBParams)
		if err != nil {
			return r, err
		}
	}

	dbIngredients, err := qtx.ListIngredientsByRecipeID(ctx, arg.ID)
	if err != nil {
		return r, err
	}

	// Update instructions in Recipe
	// List instructions from db
	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, recipe.ID)
	if err != nil {
		return r, err
	}

	// Make map from db
	dbInstructionsMap := make(map[int]Instruction, len(dbInstructions))
	for _, dbi := range dbInstructions {
		dbInstructionsMap[int(dbi.StepNo)] = dbi
	}

	var toAdd, toUpdate []step
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
		err = qtx.AddInstructionToRecipe(ctx, AddInstructionToRecipeParams{
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Instruction: pi.Instruction,
			StepNo:      int32(pi.StepNo),
			RecipeID:    recipe.ID,
		})
		if err != nil {
			return r, err
		}
	}

	for _, pi := range toUpdate {
		err = qtx.UpdateInstructionByID(ctx, UpdateInstructionByIDParams{
			UpdatedAt:   time.Now().UTC(),
			Instruction: pi.Instruction,
			StepNo:      int32(pi.StepNo),
			RecipeID:    recipe.ID,
		})
		if err != nil {
			return r, err
		}
	}

	// The rest of the dbMap is not found in param, so delete them
	for _, dbi := range dbInstructionsMap {
		err = qtx.DeleteInstructionByID(ctx, DeleteInstructionByIDParams{
			StepNo:   dbi.StepNo,
			RecipeID: dbi.RecipeID,
		})
		if err != nil {
			return r, err
		}
	}

	dbInstructions, err = qtx.ListInstructionsByRecipeID(ctx, recipe.ID)
	if err != nil {
		return r, err
	}

	// Assemble all updated data
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

	r = models.Recipe{
		ID:                recipe.ID,
		CreatedAt:         recipe.CreatedAt,
		UpdatedAt:         recipe.UpdatedAt,
		Name:              recipe.Name,
		ExternalUrl:       recipe.ExternalUrl,
		UserID:            recipe.UserID,
		Servings:          int(recipe.Servings),
		Yield:             recipe.Yield,
		CookTimeInMinutes: int(recipe.CookTimeInMinutes),
		Notes:             recipe.Notes,
		Ingredients:       ingredients,
		Instructions:      instructions,
	}

	return r, tx.Commit(ctx)
}
