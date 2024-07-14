package handlers

import (
	"net/http"
	"strconv"

	"github.com/quangd42/meal-planner/internal/models"
)

func getPaginationParamValue(r *http.Request, name string, defaultValue int32) int32 {
	val := int32(defaultValue)
	paramStr := r.URL.Query().Get(name)
	if paramStr == "" {
		return val
	}
	param64, err := strconv.ParseInt(paramStr, 10, 32)
	if err != nil {
		return val
	}
	val = int32(param64)
	return val
}

func getPaginationParams(r *http.Request) models.RecipesPagination {
	var limit, offset int32
	limit = getPaginationParamValue(r, "limit", 20)
	offset = getPaginationParamValue(r, "offset", 0)
	return models.RecipesPagination{
		Limit:  limit,
		Offset: offset,
	}
}
