package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ajg/form"
)

type Validator interface {
	Validate(ctx context.Context) error
}

func decodeJSONValidate[T Validator](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return v, err
	}

	err = v.Validate(r.Context())
	if err != nil {
		return v, err
	}

	return v, nil
}

func decodeFormValidate[T Validator](r *http.Request) (T, error) {
	var v T
	err := form.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return v, err
	}

	err = v.Validate(r.Context())
	if err != nil {
		return v, err
	}

	return v, nil
}
