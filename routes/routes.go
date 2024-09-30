package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"

	"github.com/gofiber/fiber/v2"
)

func RoutesRegister(db *platform.Postgres, api *fiber.App) {
	userRepository := repository.CreatenewUserRepository(db)
	userUsecases := usecases.CreateNewUserUsecases(userRepository)
	userController := controller.CreateNewUserController(userUsecases)

	api.Post("/register", userController.CreateUser)
}
