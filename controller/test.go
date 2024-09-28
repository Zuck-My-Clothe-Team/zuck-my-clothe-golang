package controller

import (
	_ "fmt"

	"github.com/gofiber/fiber/v2"
)

func TestController(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK).SendString("Gu Hum Yaii")
	return nil
}
