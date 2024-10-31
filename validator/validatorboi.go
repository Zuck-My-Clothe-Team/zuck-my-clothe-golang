package validatorboi

import (
	"github.com/go-playground/validator/v10"
)

// Validator instance that can be used throughout the application
var validate *validator.Validate

// Initialize function to create a new validator instance
func CreateValidator() string {
	validate = validator.New()

	if validate == nil {
		return "cannot create validator"
	}

	validate.RegisterValidation("machineType", machineTypeValidation)

	return "success"
}

func Validate(s interface{}) error {
	return validate.Struct(s)
}

// Custom validation function for MachineType
func machineTypeValidation(fl validator.FieldLevel) bool {
	machineType := fl.Field().String()

	if machineType == "Washer" || machineType == "Dryer" {
		return true
	} else {
		return false
	}
}
