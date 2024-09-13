package services

import (
	"compress/gzip"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/database"
	"github.com/quangd42/meal-planner/internal/models"
)

var ErrUnauthorized = errors.New("unauthorized")

type RecipeService struct {
	store *database.Store
}

func NewRecipeService(store *database.Store) RecipeService {
	return RecipeService{store: store}
}

func (rs RecipeService) CreateRecipe(ctx context.Context, userID uuid.UUID, arg models.RecipeRequest) (models.Recipe, error) {
	var r models.Recipe

	tx, err := rs.store.DB.Begin(ctx)
	if err != nil {
		return r, err
	}
	defer tx.Rollback(ctx)

	qtx := rs.store.Q.WithTx(tx)

	// Create host Recipe
	dbRecipe, err := qtx.CreateRecipe(ctx, database.CreateRecipeParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		Name:              arg.Name,
		ExternalUrl:       arg.ExternalURL,
		Description:       arg.Description,
		Servings:          int32(arg.Servings),
		Yield:             arg.Yield,
		CookTimeInMinutes: int32(arg.CookTimeInMinutes),
		Notes:             arg.Notes,
		UserID:            userID,
	})
	if err != nil {
		return r, err
	}

	// Add Cuisines to host Recipe
	for _, c := range arg.Cuisines {
		err = qtx.AddCuisinesToRecipe(ctx,
			database.AddCuisinesToRecipeParams{
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				RecipeID:  dbRecipe.ID,
				CuisineID: c,
			})
		if err != nil {
			return r, checkErrDBConstraint(err)
		}
	}

	dbCuisines, err := qtx.ListCuisinesByRecipeID(ctx, dbRecipe.ID)
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
			Index:        int32(p.Index),
		}
	}

	_, err = qtx.AddIngredientsToRecipe(ctx, dbIngredientParams)
	if err != nil {
		return r, checkErrDBConstraint(err)
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
			return r, checkErrDBConstraint(err)
		}
	}

	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	r = assembleWholeRecipe(dbRecipe, dbCuisines, dbIngredients, dbInstructions)

	return r, tx.Commit(ctx)
}

func (rs RecipeService) UpdateRecipeByID(ctx context.Context, userID, recipeID uuid.UUID, arg models.RecipeRequest) (models.Recipe, error) {
	var r models.Recipe

	// Check if the recipe belongs to the user
	// NOTE: this perhaps should be in a middleware
	targetRecipe, err := rs.store.Q.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return r, checkErrNoRows(err)
	}
	if targetRecipe.UserID != userID {
		return r, ErrUnauthorized
	}

	tx, err := rs.store.DB.Begin(ctx)
	if err != nil {
		return r, err
	}
	defer tx.Rollback(ctx)

	qtx := rs.store.Q.WithTx(tx)

	// Update host Recipe
	dbRecipe, err := qtx.UpdateRecipeByID(ctx, database.UpdateRecipeByIDParams{
		ID:                recipeID,
		UpdatedAt:         time.Now().UTC(),
		Name:              arg.Name,
		ExternalUrl:       arg.ExternalURL,
		Description:       arg.Description,
		Servings:          int32(arg.Servings),
		Yield:             arg.Yield,
		CookTimeInMinutes: int32(arg.CookTimeInMinutes),
		Notes:             arg.Notes,
	})
	if err != nil {
		return r, err
	}

	// Update Cuisines in Recipe
	dbCuisines, err := updateCuisinesInRecipe(ctx, qtx, arg.Cuisines, dbRecipe.ID)
	if err != nil {
		return r, checkErrDBConstraint(err)
	}

	// Update Ingredients in Recipe
	dbIngredients, err := updateIngredientsInRecipe(ctx, qtx, arg.Ingredients, dbRecipe.ID)
	if err != nil {
		return r, checkErrDBConstraint(err)
	}

	// Update instructions in Recipe
	dbInstructions, err := updateInstructionsInRecipe(ctx, qtx, arg.Instructions, dbRecipe.ID)
	if err != nil {
		return r, checkErrDBConstraint(err)
	}

	// Assemble all updated data
	r = assembleWholeRecipe(dbRecipe, dbCuisines, dbIngredients, dbInstructions)

	return r, tx.Commit(ctx)
}

func (rs RecipeService) ListRecipesByUserID(ctx context.Context, userID uuid.UUID, pgn models.RecipesPagination) ([]models.RecipeInList, error) {
	var recipes []models.RecipeInList
	dbRecipes, err := rs.store.Q.ListRecipesByUserID(ctx, database.ListRecipesByUserIDParams{
		UserID: userID,
		Limit:  pgn.Limit,
		Offset: pgn.Offset,
	})
	if err != nil {
		return recipes, err
	}

	for _, r := range dbRecipes {
		recipes = append(recipes, models.RecipeInList{
			ID:                r.ID,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			Name:              r.Name,
			ExternalURL:       r.ExternalUrl,
			ExternalImageURL:  r.ExternalImageUrl,
			Description:       r.Description,
			UserID:            r.UserID,
			Servings:          int(r.Servings),
			Yield:             r.Yield,
			CookTimeInMinutes: int(r.CookTimeInMinutes),
		})
	}

	return recipes, nil
}

func (rs RecipeService) ListRecipesWithCuisinesByUserID(ctx context.Context, userID uuid.UUID, pgn models.RecipesPagination) ([]models.RecipeInList, error) {
	var recipes []models.RecipeInList
	dbRecipes, err := rs.store.Q.ListRecipesWithCuisinesByUserID(ctx, database.ListRecipesWithCuisinesByUserIDParams{
		UserID: userID,
		Limit:  pgn.Limit,
		Offset: pgn.Offset,
	})
	if err != nil {
		return recipes, err
	}

	for _, r := range dbRecipes {
		recipes = append(recipes, models.RecipeInList{
			ID:                r.ID,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			Name:              r.Name,
			ExternalURL:       r.ExternalUrl,
			ExternalImageURL:  r.ExternalImageUrl,
			Description:       r.Description,
			UserID:            r.UserID,
			Servings:          int(r.Servings),
			Yield:             r.Yield,
			CookTimeInMinutes: int(r.CookTimeInMinutes),
			Cuisines:          string(r.Cuisines),
		})
	}

	return recipes, nil
}

func (rs RecipeService) GetRecipeByID(ctx context.Context, recipeID uuid.UUID) (models.Recipe, error) {
	var r models.Recipe

	tx, err := rs.store.DB.Begin(ctx)
	if err != nil {
		return models.Recipe{}, err
	}
	defer tx.Rollback(ctx)

	qtx := rs.store.Q.WithTx(tx)

	dbRecipe, err := qtx.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return r, checkErrNoRows(err)
	}

	dbCuisines, err := qtx.ListCuisinesByRecipeID(ctx, dbRecipe.ID)
	if err != nil {
		return r, err
	}

	dbIngredients, err := qtx.ListIngredientsByRecipeID(ctx, recipeID)
	if err != nil {
		return r, err
	}

	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, recipeID)
	if err != nil {
		return r, err
	}
	r = assembleWholeRecipe(dbRecipe, dbCuisines, dbIngredients, dbInstructions)
	return r, tx.Commit(ctx)
}

func (rs RecipeService) DeleteRecipeByID(ctx context.Context, recipeID uuid.UUID) error {
	err := rs.store.Q.DeleteRecipe(ctx, recipeID)
	if err != nil {
		return err
	}
	return nil
}

func assembleWholeRecipe(dr database.Recipe, dbCuisines []database.ListCuisinesByRecipeIDRow, dbIngredients []database.ListIngredientsByRecipeIDRow, dbInstructions []database.Instruction) models.Recipe {
	cuisines := make([]models.CuisineInRecipe, len(dbCuisines))
	for i, c := range dbCuisines {
		cuisines[i] = models.CuisineInRecipe{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	ingredients := []models.IngredientInRecipe{}
	for _, di := range dbIngredients {
		ingredients = append(ingredients, models.IngredientInRecipe{
			ID:       di.ID,
			Amount:   di.Amount,
			PrepNote: di.PrepNote,
			Name:     di.Name,
			Index:    int(di.Index),
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
		ExternalURL:       dr.ExternalUrl,
		Description:       dr.Description,
		UserID:            dr.UserID,
		Servings:          int(dr.Servings),
		Yield:             dr.Yield,
		CookTimeInMinutes: int(dr.CookTimeInMinutes),
		Notes:             dr.Notes,
		Cuisines:          cuisines,
		Ingredients:       ingredients,
		Instructions:      instructions,
	}
}

func updateCuisinesInRecipe(ctx context.Context, qtx *database.Queries, params []uuid.UUID, recipeID uuid.UUID) ([]database.ListCuisinesByRecipeIDRow, error) {
	var dbCuisines []database.ListCuisinesByRecipeIDRow
	// Add Cuisines to host Recipe
	dbCuisines, err := qtx.ListCuisinesByRecipeID(ctx, recipeID)
	if err != nil {
		return dbCuisines, err
	}
	dbCuisinesMap := make(map[uuid.UUID]bool, len(dbCuisines))
	for _, c := range dbCuisines {
		dbCuisinesMap[c.ID] = true
	}

	var toAdd, toRemove []uuid.UUID
	for _, c := range params {
		if _, ok := dbCuisinesMap[c]; !ok {
			toAdd = append(toAdd, c)
		} else {
			delete(dbCuisinesMap, c)
		}
	}
	for id := range dbCuisinesMap {
		toRemove = append(toRemove, id)
	}

	for _, c := range toAdd {
		err = qtx.AddCuisinesToRecipe(ctx,
			database.AddCuisinesToRecipeParams{
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				RecipeID:  recipeID,
				CuisineID: c,
			})
		if err != nil {
			return dbCuisines, err
		}
	}

	for _, c := range toRemove {
		err = qtx.RemoveCuisineFromRecipe(ctx, database.RemoveCuisineFromRecipeParams{
			RecipeID:  recipeID,
			CuisineID: c,
		})
		if err != nil {
			return dbCuisines, err
		}
	}

	// Get the latest list
	dbCuisines, err = qtx.ListCuisinesByRecipeID(ctx, recipeID)
	if err != nil {
		return dbCuisines, err
	}

	return dbCuisines, nil
}

func updateIngredientsInRecipe(ctx context.Context, qtx *database.Queries, params []models.IngredientInRecipe, recipeID uuid.UUID) ([]database.ListIngredientsByRecipeIDRow, error) {
	var dbIngredients []database.ListIngredientsByRecipeIDRow
	// Remove all ingredients in Recipe and add them anew
	err := qtx.RemoveAllIngredientsFromRecipe(ctx, recipeID)
	if err != nil {
		return dbIngredients, err
	}

	// Add Ingredients back to host Recipe
	dbIngredientParams := make([]database.AddIngredientsToRecipeParams, len(params))
	for i, p := range params {
		dbIngredientParams[i] = database.AddIngredientsToRecipeParams{
			RecipeID:     recipeID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			IngredientID: p.ID,
			Amount:       p.Amount,
			PrepNote:     p.PrepNote,
			Index:        int32(p.Index),
		}
	}

	_, err = qtx.AddIngredientsToRecipe(ctx, dbIngredientParams)
	if err != nil {
		return dbIngredients, err
	}

	dbIngredients, err = qtx.ListIngredientsByRecipeID(ctx, recipeID)
	if err != nil {
		return dbIngredients, err
	}

	return dbIngredients, nil
}

func updateInstructionsInRecipe(ctx context.Context, qtx *database.Queries, params []models.InstructionInRecipe, recipeID uuid.UUID) ([]database.Instruction, error) {
	var dbInstructions []database.Instruction
	// List instructions from db
	dbInstructions, err := qtx.ListInstructionsByRecipeID(ctx, recipeID)
	if err != nil {
		return dbInstructions, err
	}

	// Make map from db
	dbInstructionsMap := make(map[int]database.Instruction, len(dbInstructions))
	for _, dbi := range dbInstructions {
		dbInstructionsMap[int(dbi.StepNo)] = dbi
	}

	var toAdd, toUpdate []models.InstructionInRecipe
	for _, pi := range params {
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
			RecipeID:    recipeID,
		})
		if err != nil {
			return dbInstructions, err
		}
	}

	for _, pi := range toUpdate {
		err = qtx.UpdateInstructionByID(ctx, database.UpdateInstructionByIDParams{
			UpdatedAt:   time.Now().UTC(),
			Instruction: pi.Instruction,
			StepNo:      int32(pi.StepNo),
			RecipeID:    recipeID,
		})
		if err != nil {
			return dbInstructions, err
		}
	}

	// The rest of the dbMap is not found in param, so delete them
	for _, dbi := range dbInstructionsMap {
		err = qtx.DeleteInstructionByID(ctx, database.DeleteInstructionByIDParams{
			StepNo:   dbi.StepNo,
			RecipeID: dbi.RecipeID,
		})
		if err != nil {
			return dbInstructions, err
		}
	}

	dbInstructions, err = qtx.ListInstructionsByRecipeID(ctx, recipeID)
	if err != nil {
		return dbInstructions, err
	}

	return dbInstructions, nil
}

func (rs RecipeService) SaveExternalImage(recipeID uuid.UUID, url *string) {
	ctx := context.Background()
	if url == nil || *url == "" {
		return
	}
	imageURL, err := fetchOGImage(*url)
	if err != nil {
		log.Println(err, url)
	}
	err = rs.store.Q.SaveExternalImageURL(ctx, database.SaveExternalImageURLParams{
		ID:               recipeID,
		ExternalImageUrl: &imageURL,
	})
	if err != nil {
		log.Println(err, url)
	}
}

func fetchOGImage(url string) (string, error) {
	og := opengraph.NewOpenGraph()

	agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:103.0) Gecko/20100101 Firefox/103.0"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	// req.Header.Set("User-Agent", agent)
	// req.Header.Set("Accept-Encoding", "gzip, deflate")
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Connection", "keep-alive")
	req.Header = http.Header{
		"User-Agent":      {agent},
		"Accept-Encoding": {"gzip, deflate"},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
	}

	// make the http request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err, url)
	}
	defer resp.Body.Close()

	// decompress the response
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Println(err, url)
	}
	defer reader.Close()

	err = og.ProcessHTML(reader)
	if err != nil {
		return "", errors.New("html cannot be processed")
	}

	if len(og.Images) == 0 {
		return "", errors.New("no OG image")
	}

	return og.Images[0].URL, nil
}
