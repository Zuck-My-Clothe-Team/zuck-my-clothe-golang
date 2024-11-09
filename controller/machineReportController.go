package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type MachineReportController interface {
	CreateMachineReport(c *fiber.Ctx) error
	FindMachineReportByUserID(c *fiber.Ctx) error
	FindMachineReportByBranch(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	UpdateMachineReportStatus(c *fiber.Ctx) error
	DeleteMachineReport(c *fiber.Ctx) error
}

type machineReportController struct {
	machineReportUsecase model.MachineReportsUsecase
}

// tokenClaimer extracts the user ID or other claims from the JWT token
func tokenClaimer(c *fiber.Ctx, localParam string, field string) string {
	token := c.Locals(localParam).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	result := claims[field].(string)
	return result
}

// CreateNewMachineReportController creates a new instance of the machineReportController
func CreateNewMachineReportController(machineReportUsecase model.MachineReportsUsecase) MachineReportController {
	return &machineReportController{machineReportUsecase: machineReportUsecase}
}

//	@Summary		Create a new machine report
//	@Description	Create a new machine report for a specific user
//	@Tags			Machine Reports
//	@Accept			json
//	@Produce		json
//	@Param			report	body		model.AddMachineReportDTO	true	"Machine Report"
//	@Success		200		{object}	model.MachineReportDetail
//	@Failure		406		{string}	string	"Invalid request"
//	@Failure		202		{string}	string	"Validation error"
//	@Security		BearerAuth
//	@Router			/report/add [post]
func (u *machineReportController) CreateMachineReport(c *fiber.Ctx) error {
	userID := tokenClaimer(c, "user", "userID")
	newReport := new(model.AddMachineReportDTO)
	if err := c.BodyParser(newReport); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	if err := validatorboi.Validate(newReport); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	createdReport, createErr := u.machineReportUsecase.CreateMachineReport(newReport, userID)
	if createErr != nil {
		return c.Status(fiber.StatusAccepted).SendString(createErr.Error())
	}
	return c.Status(fiber.StatusOK).JSON(createdReport)
}

//	@Summary		Get machine reports by user ID
//	@Description	Retrieve all machine reports associated with a specific user
//	@Tags			Machine Reports
//	@Produce		json
//	@Success		200	{array}		model.MachineReportDetail
//	@Failure		202	{string}	string	"Accepted"
//	@Security		BearerAuth
//	@Router			/report/user/{userID} [get]
func (u *machineReportController) FindMachineReportByUserID(c *fiber.Ctx) error {
	userID := tokenClaimer(c, "user", "userID")
	reportList, err := u.machineReportUsecase.FindMachineReportByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(reportList)
}

//	@Summary		Get machine reports by branch ID
//	@Description	Retrieve all machine reports for a specific branch
//	@Tags			Machine Reports
//	@Produce		json
//	@Param			branchID	path		string	true	"Branch ID"
//	@Success		200			{array}		model.MachineReportDetail
//	@Failure		202			{string}	string	"No records found"
//	@Failure		204			{string}	string	"Invalid request"
//	@Security		BearerAuth
//	@Router			/report/branch/{branchID} [get]
func (u *machineReportController) FindMachineReportByBranch(c *fiber.Ctx) error {
	userID := tokenClaimer(c, "user", "userID")
	branchID := c.Params("branchID")
	result, err := u.machineReportUsecase.FindMachineReportByBranch(branchID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		}
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Get all machine reports
//	@Description	Retrieve all machine reports in the system
//	@Tags			Machine Reports
//	@Produce		json
//	@Success		200	{array}		model.MachineReports
//	@Failure		500	{string}	string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/report/ [get]
func (u *machineReportController) GetAll(c *fiber.Ctx) error {
	reportList, err := u.machineReportUsecase.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(reportList)
}

//	@Summary		Update the status of a machine report
//	@Description	Update the status of a specific machine report
//	@Tags			Machine Reports
//	@Accept			json
//	@Produce		json
//	@Param			status	body		model.UpdateMachineReportStatusDTO	true	"Updated Status"
//	@Success		200		{object}	model.MachineReportDetail
//	@Failure		204		{string}	string	"Record not found"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Invalid request"
//	@Security		BearerAuth
//	@Router			/report/update [put]
func (u *machineReportController) UpdateMachineReportStatus(c *fiber.Ctx) error {
	updateReport := new(model.UpdateMachineReportStatusDTO)
	userID := tokenClaimer(c, "user", "userID")
	userRole := tokenClaimer(c, "user", "positionID")
	if err := c.BodyParser(updateReport); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	updatedReport, err := u.machineReportUsecase.UpdateMachineReportStatus(*updateReport, userID, userRole)
	if err != nil {
		if err.Error() == "un authorized" {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		} else if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString("record not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(updatedReport)
}

//	@Summary		Delete a machine report
//	@Description	Delete a specific machine report by ID
//	@Tags			Machine Reports
//	@Param			reportID	path		string	true	"Report ID"
//	@Success		200			{string}	string	"Successfully deleted"
//	@Failure		204			{string}	string	"Record not found"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Invalid request"
//	@Security		BearerAuth
//	@Router			/report/delete/{reportID} [delete]
func (u *machineReportController) DeleteMachineReport(c *fiber.Ctx) error {
	reportID := c.Params("reportID")
	userID := tokenClaimer(c, "user", "userID")
	userRole := tokenClaimer(c, "user", "positionID")
	err := u.machineReportUsecase.DeleteMachineReport(reportID, userID, userRole)
	if err != nil {
		if err.Error() == "un authorized" {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		} else if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString("record not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString("Successfully deleted")
}
