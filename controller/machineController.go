package controller

import (
	"strconv"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MachineController interface {
	AddMachine(c *fiber.Ctx) error
	GetByBranchID(c *fiber.Ctx) error
	GetByMachineSerial(c *fiber.Ctx) error
	GetAvailableMachineInBranch(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	UpdateLabel(c *fiber.Ctx) error
	UpdateActive(c *fiber.Ctx) error
	SoftDelete(c *fiber.Ctx) error
}

type machineController struct {
	machineUsecase usecases.MachineUsecase
}

func CreateMachineController(machineUsecase usecases.MachineUsecase) MachineController {
	return &machineController{machineUsecase: machineUsecase}
}

// @Summary		Add new machine
// @Description	Add a new machine to the system
// @Tags			Machine
// @Accept			json
// @Produce		json
// @Param			MachineModel	body		model.AddMachine	true	"New Machine Data"
// @Success		201				{object}	model.Machine		"Created"
// @Failure		406				{string}	string				"Not Acceptable"
// @Failure		500				{string}	string				"internal server error"
// @Router			/machine/add [post]
func (u *machineController) AddMachine(c *fiber.Ctx) error {
	new_machine := new(model.AddMachine)

	if err := c.BodyParser(new_machine); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if err := validatorboi.Validate(new_machine); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	new_machine.CreatedBy = claims["userID"].(string)

	response, err := u.machineUsecase.AddMachine(new_machine)

	if err != nil {
		if err.Error() == "null detected on one or more essential field(s)" {
			return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// @Summary		Get machine details by serial
// @Description	Get details of a specific machine by its serial number
// @Tags			Machine
// @Produce		json
// @Param			serial_id	path		string			true	"Machine Serial ID"
// @Success		200			{object}	model.Machine	"OK"
// @Failure		404			{string}	string			"Not Found"
// @Failure		500			{string}	string			"Internal Server Error"
// @Router			/machine/detail/{serial_id} [get]
func (u *machineController) GetByMachineSerial(c *fiber.Ctx) error {
	serialID := c.Params("serial_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["positionID"].(string)

	var isAdminView bool = false

	if role == "BranchManager" || role == "SuperAdmin" {
		isAdminView = true
	}

	branch, err := u.machineUsecase.GetByMachineSerial(serialID, isAdminView)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {

			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.Status(fiber.StatusOK).JSON(branch)
}

// @Summary		Get available machines by branch ID
// @Description	Get all available machines under a specific branch
// @Tags			Machine
// @Produce		json
// @Param			branch_id	path		string			true	"Branch ID"
// @Success		200			{object}	model.Machine	"OK"
// @Failure		404			{string}	string			"Not Found"
// @Failure		202			{string}	string			"Accepted"
// @Router			/machine/available/branch/{branch_id} [get]
func (u *machineController) GetAvailableMachineInBranch(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")
	response, err := u.machineUsecase.GetAvailableMachineInBranch(branchID)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary		Get machines by branch ID
// @Description	Get all machines under a specific branch
// @Tags			Machine
// @Produce		json
// @Param			branch_id	path		string			true	"Branch ID"
// @Success		200			{object}	model.Machine	"OK"
// @Failure		404			{string}	string			"Not Found"
// @Failure		500			{string}	string			"Internal Server Error"
// @Router			/machine/branch/{branch_id} [get]
func (u *machineController) GetByBranchID(c *fiber.Ctx) error {
	branch_id := c.Params("branch_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["positionID"].(string)

	var isAdminView bool = false

	if role == "BranchManager" || role == "SuperAdmin" {
		isAdminView = true
	}

	result, err := u.machineUsecase.GetByBranchID(branch_id, isAdminView)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {

			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary		Get all machines
// @Description	Retrieve all machines in the system
// @Tags			Machine
// @Produce		json
// @Success		200	{array}		model.Machine	"OK"
// @Failure		404	{string}	string			"Not Found"
// @Failure		500	{string}	string			"Internal Server Error"
// @Router			/machine/all [get]
func (u *machineController) GetAll(c *fiber.Ctx) error {
	result, err := u.machineUsecase.GetAll()

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary		Soft delete machine
// @Description	Soft delete a machine by its serial ID
// @Tags			Machine
// @Param			serial_id	path		string			true	"Machine Serial ID"
// @Success		200			{object}	model.Machine	"OK"
// @Failure		404			{string}	string			"Not Found"
// @Failure		500			{string}	string			"Internal Server Error"
// @Router			/machine/delete/{serial_id} [delete]
func (u *machineController) SoftDelete(c *fiber.Ctx) error {
	serial_id := c.Params("serial_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	deleted_by := claims["userID"].(string)

	result, err := u.machineUsecase.SoftDelete(serial_id, deleted_by)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary		Update machine active status
// @Description	Set the active status of a machine
// @Tags			Machine
// @Param			serial_id	path		string			true	"Machine Serial ID"
// @Param			set_active	path		string			true	"Set Active (true/false)"
// @Success		200			{object}	model.Machine	"OK"
// @Failure		400			{string}	string			"Bad Request"
// @Failure		404			{string}	string			"Not Found"
// @Failure		500			{string}	string			"Internal Server Error"
// @Router			/machine/update/{serial_id}/set_active/{set_active} [put]
func (u *machineController) UpdateActive(c *fiber.Ctx) error {
	machine_serial := c.Params("serial_id")
	set_active_param := c.Params("set_active")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	updated_by := claims["userID"].(string)

	var set_active bool

	if set_active_param == "true" {
		set_active = true
	} else if set_active_param == "false" {
		set_active = false
	} else {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	result, err := u.machineUsecase.UpdateActive(machine_serial, set_active, updated_by)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary		Update machine label
// @Description	Update machine label
// @Tags			Machine
// @Param			serial_id	path		string			true	"Machine Serial ID"
// @Param			label		path		string			true	"New label (int)"
// @Success		200			{object}	model.Machine	"OK"
// @Success		204			{string}	string			"Not Content"
// @Failure		404			{string}	string			"Not Found"
// @Failure		500			{string}	string			"Internal Server Error"
// @Router			/machine/update/{serial_id}/set_label/{label} [put]
func (u *machineController) UpdateLabel(c *fiber.Ctx) error {
	machine_serial := c.Params("serial_id")
	label := c.Params("label")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	updated_by := claims["userID"].(string)

	intLabel, err := strconv.Atoi(label)

	if err != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	result, err := u.machineUsecase.UpdateLabel(machine_serial, intLabel, updated_by)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
