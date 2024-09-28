package controller

import (
	_ "fmt"

	"github.com/gofiber/fiber/v2"
)

// @Summary		Test Controller
// @Description	Returns a simple message to verify the endpoint.
// @Tags			Test
// @Accept			json
// @Produce		plain
// @Success		200	{string}	string	"Gu Hum Yaii"
// @Router			/testrout [get]
func TestController(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK).SendString("Gu Hum Yaii")
	return nil
}
