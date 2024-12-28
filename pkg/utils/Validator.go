package utils

import "github.com/go-playground/validator/v10"

type ValidationErrors []string

// ValidateStruct validates the struct based on the tags defined in the struct
func ValidateStruct(data interface{}) (ValidationErrors, bool) {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		validationErrors := ValidationErrors{}
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, e.Field()+" is invalid: ")
		}
		return validationErrors, false
	}
	return nil, true
}
