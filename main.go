package main

import (
	"log"
	"zuck-my-clothe/zuck-my-clothe-backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	_, err := config.Load()
	if err != nil {
		log.Fatal("Can't load config", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	api := fiber.New()
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8081",
		AllowCredentials: true,
	}))

	api.Listen(":3000")
}
