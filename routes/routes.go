package routes

import (
	"zuck-my-clothe/zuck-my-clothe-backend/controller"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"github.com/gofiber/fiber/v2"
)

func RoutesRegister(db *platform.Postgres, api *fiber.App) {
	api.Get("/testrout", controller.TestController)
}
