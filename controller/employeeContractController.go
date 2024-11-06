package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type EmployeeContractController interface {
	CreateEmployeeContract(c *fiber.Ctx) error
	SoftDelete(c *fiber.Ctx) error
	GetByBranchID(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type employeeContractController struct {
	usecase usecases.EmployeeContractUsecases
}

func CreateNewEmployeeContractController(employeeContractUsecases usecases.EmployeeContractUsecases) EmployeeContractController {
	return &employeeContractController{usecase: employeeContractUsecases}
}

// CreateEmployeeContract godoc
//
//	@Summary		Create a new employee contract
//	@Description	Create a new employee contract
//	@Tags			EmployeeContract
//	@Accept			json
//	@Produce		json
//	@Param			contract	body		model.EmployeeContractDTO	true	"Employee Contract"
//	@Success		201			{string}	string						"Created"
//	@Failure		406			{string}	string						"Not Acceptable"
//	@Failure		403			{string}	string						"Forbidden"
//	@Router			/employee-contract [post]
func (controller *employeeContractController) CreateEmployeeContract(c *fiber.Ctx) error {
	newContract := new(model.EmployeeContractDTO)
	if err := c.BodyParser(newContract); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	if err := controller.usecase.CreateEmployeeContract(newContract); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

// SoftDelete godoc
//
//	@Summary		Soft delete an employee contract
//	@Description	Soft delete an employee contract
//	@Tags			EmployeeContract
//	@Param			contract_id	path		string	true	"Contract ID"
//	@Success		200			{object}	model.EmployeeContract
//	@Failure		204			{string}	string	"No Content"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/employee-contract/{contract_id} [delete]
func (controller *employeeContractController) SoftDelete(c *fiber.Ctx) error {
	contractId := c.Params("contract_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userID"].(string)

	deletedUser, err := controller.usecase.SoftDelete(contractId, userId)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(deletedUser)
}

// GetByBranchID godoc
//
//	@Summary		Get employee contracts by branch ID
//	@Description	Get employee contracts by branch ID
//	@Tags			EmployeeContract
//	@Param			branch_id	path		string	true	"Branch ID"
//	@Success		200			{array}		model.EmployeeContract
//	@Failure		204			{string}	string	"No Content"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/employee-contract/branch/{branch_id} [get]
func (controller *employeeContractController) GetByBranchID(c *fiber.Ctx) error {
	branchId := c.Params("branch_id")

	contracts, err := controller.usecase.GetByBranchID(branchId)
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(contracts)
}

// GetByUserID godoc
//
//	@Summary		Get employee contracts by user ID
//	@Description	Get employee contracts by user ID
//	@Tags			EmployeeContract
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{array}		model.EmployeeContract
//	@Failure		204		{string}	string	"No Content"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/employee-contract/user/{user_id} [get]
func (controller *employeeContractController) GetByUserID(c *fiber.Ctx) error {
	userId := c.Params("user_id")
	contracts, err := controller.usecase.GetByUserID(userId)
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(contracts)
}

// GetAll godoc
//
//	@Summary		Get all employee contracts
//	@Description	Get all employee contracts
//	@Tags			EmployeeContract
//	@Success		200	{array}		model.EmployeeContract
//	@Failure		204	{string}	string	"No Content"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/employee-contract [get]
func (controller *employeeContractController) GetAll(c *fiber.Ctx) error {
	contracts, err := controller.usecase.GetAll()
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(contracts)
}
