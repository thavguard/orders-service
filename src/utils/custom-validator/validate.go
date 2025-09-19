package customvalidator

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func zipUS(fl validator.FieldLevel) bool {

	// Регулярка для US ZIP
	zipUSRegex := regexp.MustCompile(`^\d{5}(-\d{4})?$`)

	val := fl.Field().String()

	if val == "" {
		return false
	}

	return zipUSRegex.MatchString(val)
}

func NewValidator() (*validator.Validate, error) {

	validate := validator.New(validator.WithRequiredStructEnabled()) // рекомендовано

	if err := validate.RegisterValidation("zipcode", zipUS); err != nil {
		log.Printf("ERROR IN validate.RegisterValidation: %v\n", err)
		return &validator.Validate{}, err
	}

	return validate, nil
}
