package validatorboi

import (
	"reflect"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

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
	validate.RegisterValidation("serviceType", serviceTypeValidation)
	validate.RegisterValidation("orderStatus", orderStatusValidation)
	validate.RegisterValidation("requiredBool", requiredBool)
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

func requiredBool(fl validator.FieldLevel) bool {
	field := fl.Field()
	return field.Kind() == reflect.Bool
}

func serviceTypeValidation(fl validator.FieldLevel) bool {
	serviceType := fl.Field().String()

	if serviceType == "Washing" ||
		serviceType == "Drying" ||
		serviceType == "Delivery" {
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


func orderStatusValidation(fl validator.FieldLevel) bool {
	orderStatus := fl.Field().String()

	if orderStatus == "Waiting" ||
		orderStatus == "Processing" ||
		orderStatus == "Completed" ||
		orderStatus == "Canceled" {
		return true
	} else {
		return false
	}
}

func userRolesValidation(fl validator.FieldLevel) bool {
	userRoles := fl.Field().String()

	if userRoles == string(model.SuperAdmin) || userRoles == string(model.BranchManager) || userRoles == string(model.Employee) || userRoles == string(model.Client) {
		return true
	} else {
		return false
	}
}
