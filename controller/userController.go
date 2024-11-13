package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/middleware"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController interface {
	CreateUser(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	GetBranchEmployee(c *fiber.Ctx) error
	GetAllManager(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	UpdateUserPassword(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	DeleteEmployeeFromBranch(c *fiber.Ctx) error
}

type userController struct {
	employeeContractUsecase usecases.EmployeeContractUsecases
	branchUsecase           usecases.BranchUsecase
	usecase                 usecases.UserUsecases
	config                  *config.Config
}

func CreateNewUserController(usecase usecases.UserUsecases, config *config.Config, employeeContractUsecase usecases.EmployeeContractUsecases, branchUsecase usecases.BranchUsecase) UserController {
	return &userController{
		usecase:                 usecase,
		config:                  config,
		employeeContractUsecase: employeeContractUsecase,
		branchUsecase:           branchUsecase,
	}
}

// @Summary		Create new user
// @Description	Create a new user by using User model
// @Tags			Users
// @Produce		json
// @Accept			json
// @Param			UserModel	body	model.UserDTO	true	"New User Data"
// @Success		201
// @Failure		403	{string}	string	"Forbidden"
// @Failure		406	{string}	string	"Not Acceptable"
// @Router			/users/ [POST]
func (controller *userController) CreateUser(c *fiber.Ctx) error {

	body := new(model.UserDTO)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	newUser := &model.Users{
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		Role:            body.Role,
		Phone:           body.Phone,
		ProfileImageURL: body.ProfileImageURL,
	}
	userdata, err := controller.usecase.CreateUser(*newUser)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	if newUser.Role == model.Employee {
		middleware.AuthRequire(c)
		tokenInterface := c.Locals("user")
		if tokenInterface == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		} else {
			token := tokenInterface.(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			userId := claims["userID"].(string)
			role := claims["positionID"].(string)

			if role != string(model.BranchManager) {
				return c.Status(fiber.StatusForbidden).SendString("You are not allowed to create this user")
			}

			ownBranches, err := controller.branchUsecase.GetByBranchOwner(userId)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			contracts := &body.Contracts

			if len(*contracts) == 0 {
				return c.Status(fiber.StatusNotAcceptable).SendString("Employee contract is required")
			} else {
				for _, contract := range *contracts {
					if contract.PositionId != string(model.Worker) && contract.PositionId != string(model.Deliver) {
						return c.Status(fiber.StatusNotAcceptable).SendString("Invalid position")
					}
					if contract.BranchID == "" {
						return c.Status(fiber.StatusNotAcceptable).SendString("Branch ID is required")
					}
				}

				validBranch := false
				for _, contract := range *contracts {
					validBranch = false
					for _, branch := range *ownBranches {
						if branch.BranchID == contract.BranchID {
							validBranch = true
							break
						}
					}
					if !validBranch {
						return c.Status(fiber.StatusNotAcceptable).SendString("Invalid branch")
					}
				}

				for _, contract := range *contracts {
					newContractDTO := &model.EmployeeContractDTO{
						UserID:     userdata.UserID,
						BranchID:   contract.BranchID,
						PositionId: contract.PositionId,
						CreatedBy:  userId,
					}
					if err := controller.employeeContractUsecase.CreateEmployeeContract(newContractDTO); err != nil {
						return c.Status(fiber.StatusForbidden).SendString(err.Error())
					}
				}
			}
		}
	}
	return c.SendStatus(fiber.StatusCreated)

}

// @Summary		Get all users
// @Description	Retrieve all users from the database
// @Tags			Users
// @Produce		json
// @Success		200	{array}		model.UserBranch[]
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/users/all [get]
func (controller *userController) GetAll(c *fiber.Ctx) error {
	users, err := controller.usecase.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(users)
}

// @Summary		Get a user by ID
// @Description	Get a user by ID from the database
// @Tags			Users
// @Param			id	path		string	true	"User ID"
// @Success		200	{object}	model.Users
// @Failure		204	{string}	string	"record not found"
// @Failure		500	{string}	string	"internal server error"
// @Router			/users/{id} [get]
func (controller *userController) GetUserById(c *fiber.Ctx) error {
	userID := c.Params("id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUserID := claims["userID"].(string)
	role := claims["positionID"].(string)

	if role == string(model.Client) && userID != tokenUserID {
		return c.Status(fiber.StatusForbidden).SendString("You are not allowed to fetch this data")
	}

	user, err := controller.usecase.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Get employees by branch ID
// @Description	Get employees by branch ID
// @Tags			Users
// @Param			branch_id	path		string	true	"Branch ID"
// @Success		200			{array}		model.UserContract
// @Failure		204			{string}	string	"No Content"
// @Failure		500			{string}	string	"Internal Server Error"
// @Router			/users/branch/{branch_id} [get]
func (controller *userController) GetBranchEmployee(c *fiber.Ctx) error {
	branchId := c.Params("branch_id")
	users, err := controller.usecase.GetBranchEmployee(branchId)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.JSON(users)
}

// @Summary		Get all managers
// @Description	Get a list of all managers
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{array}		model.Users[]
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/users/manager/all [get]
func (controller *userController) GetAllManager(c *fiber.Ctx) error {
	user, err := controller.usecase.GetAllManager()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(user)
}

// @Summary		Delete user
// @Description	Delete a user by ID
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"User ID"
// @Success		200	{struct}	model.Users
// @Failure		500	{string}	string	"Internal Server Error"
// @Router			/users/:id [delete]
func (controller *userController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	deletedUser, err := controller.usecase.DeleteUser(userID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}
	return c.Status(fiber.StatusOK).JSON(deletedUser)
}

func (controller *userController) DeleteEmployeeFromBranch(c *fiber.Ctx) error {
	userID := c.Params("id")
	branchID := c.Params("branch_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUserID := claims["userID"].(string)
	role := claims["positionID"].(string)

	if role != string(model.BranchManager) {
		return c.Status(fiber.StatusForbidden).SendString("You are not allowed to delete this user")
	}

	branches, err := controller.branchUsecase.GetByBranchOwner(tokenUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	validBranch := false
	for _, branch := range *branches {
		if branch.BranchID == branchID {
			validBranch = true
			break
		}
	}

	if !validBranch {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid branch")
	}

	contracts, err := controller.employeeContractUsecase.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(*contracts) == 0 {
		return c.Status(fiber.StatusNoContent).SendString("record not found")
	}

	for _, contract := range *contracts {
		if contract.BranchID == branchID {
			_, err := controller.employeeContractUsecase.SoftDelete(contract.ContractId, tokenUserID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
	}

	deletedUser, err := controller.usecase.DeleteUser(userID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(deletedUser)

}

// @Summary		Update user
// @Description	Update user details by ID
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			id		path		string				true	"User ID"
// @Param			User	body		model.UserUpdateDTO	true	"Updated User Data"
// @Success		200		{object}	model.Users
// @Failure		403		{string}	string	"Forbidden"
// @Failure		406		{string}	string	"Not Acceptable"
// @Failure		500		{string}	string	"Internal Server Error"
// @Router			/users/{id} [patch]
func (contr *userController) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUserID := claims["userID"].(string)
	role := claims["positionID"].(string)

	body := new(model.UserUpdateDTO)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if role != string(model.SuperAdmin) && userID != tokenUserID {
		return c.Status(fiber.StatusForbidden).SendString("You are not allowed to update")
	}

	err := contr.usecase.UpdateUser(userID, *body, role)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)

}

// UpdateUserPassword handles the request to update a user's password.
//
//	@Summary		Update user password
//	@Description	Update the password of a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			body	body		model.UserUpdatePasswordDTO	true	"New password data"
//	@Success		200		{string}	string						"Password updated successfully"
//	@Failure		403		{string}	string						"You are not allowed to update"
//	@Failure		406		{string}	string						"Invalid request body"
//	@Failure		500		{string}	string						"Internal server error"
//	@Router			/users/{id}/password [patch]
func (contr *userController) UpdateUserPassword(c *fiber.Ctx) error {
	userID := c.Params("id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUserID := claims["userID"].(string)
	role := claims["positionID"].(string)

	body := new(model.UserUpdatePasswordDTO)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if role != string(model.SuperAdmin) && userID != tokenUserID {
		return c.Status(fiber.StatusForbidden).SendString("You are not allowed to update")
	}

	err := contr.usecase.UpdateUserPassword(userID, *body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
