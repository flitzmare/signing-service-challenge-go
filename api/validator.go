package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// validateRequest validates a struct and returns formatted error messages
func validateRequest(s interface{}) []string {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required", err.Field()))
			case "oneof":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", err.Field()))
			}
		}
		return validationErrors
	}
	return nil
}