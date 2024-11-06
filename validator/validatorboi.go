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
	validate.RegisterValidation("employeeContractPosition", employeeContractValidation)
	validate.RegisterValidation("userRoles", userRolesValidation)

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

func employeeContractValidation(fl validator.FieldLevel) bool {
	employeeContract := fl.Field().String()

	if employeeContract == "Worker" || employeeContract == "Deliver" {
		return true
	} else {
		return false
	}
}

func userRolesValidation(fl validator.FieldLevel) bool {
	userRoles := fl.Field().String()

	if userRoles == "Superadmin" || userRoles == "Manager" || userRoles == "Employee" || userRoles == "Client" {
		return true
	} else {
		return false
	}
}
