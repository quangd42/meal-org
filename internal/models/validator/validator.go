package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type ValidationErrors map[string][]string

func NewValidationErrors() ValidationErrors {
	return make(ValidationErrors)
}

// Error is intended for use in development + debugging and not intended to be a production error message.
// It allows ValidationErrors to subscribe to the Error interface.
func (e ValidationErrors) Error() string {
	indent, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return string(indent)
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		valErrs := NewValidationErrors()
		for _, err := range err.(validator.ValidationErrors) {
			var msg string
			switch err.Tag() {
			case "email":
				msg = "invalid email format"
			case "min":
				if err.Param() == "0" {
					msg = "cannot be empty"
					break
				}
				msg = fmt.Sprintf("must be at least %s character long", err.Param())
			default:
				msg = err.Tag()
			}
			valErrs[err.Field()] = append(valErrs[err.Field()], msg)
		}

		return valErrs
	}
	return err
}
