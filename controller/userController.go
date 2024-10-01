package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	CreateUser(c *fiber.Ctx) error
}

type userController struct {
	usecase model.UserUsecases
}

func CreateNewUserController(usecase model.UserUsecases) UserController {
	return &userController{usecase: usecase}
}

func (controller *userController) CreateUser(c *fiber.Ctx) error {
	newUser := new(model.Users)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	if err := controller.usecase.CreateUser(*newUser); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
