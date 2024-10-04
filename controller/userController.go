package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	CreateUser(c *fiber.Ctx) error
}

type userController struct {
	usecase model.UserUsecases
	config  *config.Config
}

func CreateNewUserController(usecase model.UserUsecases, config *config.Config) UserController {
	return &userController{usecase: usecase, config: config}
}

// @Summary		Create new user
// @Description	Create a new user by using User model
// @Tags			User Controller
// @Produce		json
// @Accept			json
// @Param			UserModel	body	model.Users	true	"New User Data"
// @Success		201
// @Failure		403	{string}	string	"Forbidden"
// @Failure		406	{string}	string	"Not Acceptable"
// @Router			/register [POST]
func (controller *userController) CreateUser(c *fiber.Ctx) error {
	newUser := new(model.Users)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	newUser.GoogleID = ""
	if err := controller.usecase.CreateUser(*newUser); err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}
