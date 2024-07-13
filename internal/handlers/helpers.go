package handlers

import (
	"net/http"
	"strconv"
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
