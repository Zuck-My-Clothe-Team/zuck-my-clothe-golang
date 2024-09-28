package main

import (
	"log"
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
	"zuck-my-clothe/zuck-my-clothe-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	config, err := config.Load()
	if err != nil {
		log.Fatal("Can't load config", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	db, dbErr := platform.InitDB(config.DB_DSN)

	if dbErr != nil {
		log.Fatal("Can not Init Database", dbErr)
	}

	api := fiber.New()
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8081",
		AllowCredentials: true,
	}))

	routes.RoutesRegister(db, api)

	api.Listen(":" + config.PORT)
}
