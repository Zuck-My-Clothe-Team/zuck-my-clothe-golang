package middleware

import (
	"log"
	"strings"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleWare(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

func AuthRequire(c *fiber.Ctx) error {

	// Try Getting from Cookies first
	// If request doesn't have cookie try getting from https header
	// If header also not exist then return with unauthorized status
	reqToken := c.Cookies("jwt")
	if utils.CheckStraoPling(reqToken) {
		reqToken = string(c.Request().Header.Peek("Authorization"))
		if utils.CheckStraoPling(reqToken) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		splitToken := strings.TrimPrefix(string(reqToken), "Bearer ")
		reqToken = splitToken
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(reqToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatal("Cannot load config", err)
			}
			return []byte(cfg.JWT_ACCESS_TOKEN), nil
		})
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//Token Expiration Check
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if exp.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).SendString("token expired")
	}

	c.Locals("user", token)
	return c.Next()
}

func Claimer(c *fiber.Ctx) jwt.MapClaims {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return claims
}

func IsSuperAdmin(c *fiber.Ctx) error {
	claims := Claimer(c)
	if claims["positionID"] != "SuperAdmin" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

func IsBranchManager(c *fiber.Ctx) error {
	claims := Claimer(c)
	if claims["positionID"] != "SuperAdmin" && claims["positionID"] != "BranchManager" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

func IsEmployee(c *fiber.Ctx) error {
	claims := Claimer(c)
	if claims["positionID"] != "SuperAdmin" && claims["positionID"] != "BranchManager" && claims["positionID"] != "Employee" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
