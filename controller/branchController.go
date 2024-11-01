package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BranchController interface {
	CreateBranch(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByBranchID(c *fiber.Ctx) error
	GetByBranchOwner(c *fiber.Ctx) error
	UpdateBranch(c *fiber.Ctx) error
	DeleteBranch(c *fiber.Ctx) error
}

type branchController struct {
	branchUsecase model.BranchUsecase
}

func CreateNewBranchController(branchUsecase model.BranchUsecase) BranchController {
	return &branchController{branchUsecase: branchUsecase}
}

// @Summary		Create new branch
// @Description	Create a new branch by using Branch model
// @Tags			Branches
// @Produce		json
// @Accept			json
// @Param			BranchModel	body	model.Branch	true	"New Branch Data"
// @Success		201
// @Failure		403	{string}	string	"Forbidden"
// @Failure		406	{string}	string	"Not Acceptable"
// @Router			/branch/create [POST]
func (u *branchController) CreateBranch(c *fiber.Ctx) error {
	newBranch := new(model.Branch)
	if err := c.BodyParser(newBranch); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	claims := middleware.Claimer(c)
	if claims["userID"] == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID := claims["userID"].(string)

	response, err := u.branchUsecase.CreateBranch(newBranch, userID)
	if err != nil {
		if err.Error() == "null detected on one or more essential field(s)" {
			return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (u *branchController) GetAll(c *fiber.Ctx) error {
	branchList, err := u.branchUsecase.GetAll()
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(branchList)
}

// @Summary		Get a branch by ID
// @Description	Retrieve a single branch from the database based on its ID.
// @Tags			Branches
// @Produce		json
// @Param			id	path		string	true	"branch ID"
// @Success		200	{object}	model.Branch
// @Failure		404	{string}	string	"Not Found"
// @Router			/branch/{id} [GET]
func (u *branchController) GetByBranchID(c *fiber.Ctx) error {
	branchID := c.Params("id")
	branch, err := u.branchUsecase.GetByBranchID(branchID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(branch)
}

// @Summary		Get branch by owner
// @Description	Get branch details by branch owner
// @Tags			Branches
// @Accept			json
// @Produce		json
// @Success		200	{object}	model.Branch
// @Failure		404	{string}	string	"record not found"
// @Failure		500	{string}	string	"internal server error"
// @Router			/branches/owner [get]
func (u *branchController) GetByBranchOwner(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	branch, err := u.branchUsecase.GetByBranchOwner(claims["userID"].(string))
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {

			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.Status(fiber.StatusOK).JSON(branch)
}

// @Summary		Update branch
// @Description	Update branch details
// @Tags			Branches
// @Accept			json
// @Produce		json
// @Param			branch	body		model.Branch	true	"Branch data"
// @Success		200		{object}	model.Branch
// @Failure		406		{string}	string	"not acceptable"
// @Router			/branches [put]
func (u *branchController) UpdateBranch(c *fiber.Ctx) error {
	branch := new(model.Branch)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if err := c.BodyParser(branch); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	response, err := u.branchUsecase.UpdateBranch(branch, claims["positionID"].(string))

	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary		Delete branch
// @Description	Delete a branch
// @Tags			Branches
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Branch ID"
// @Success		200	{string}	string	"ok"
// @Failure		404	{string}	string	"record not found"
// @Failure		500	{string}	string	"internal server error"
// @Router			/branches/{id} [delete]
func (u *branchController) DeleteBranch(c *fiber.Ctx) error {
	branch := new(model.Branch)
	branchID := c.Params("id")
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	branch.BranchID = branchID
	branch.DeletedBy = claims["userID"].(string)

	err := u.branchUsecase.DeleteBranch(branch)
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.SendStatus(fiber.StatusOK)
}
