package main

import (
	"fmt"
	"log"
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	_ "zuck-my-clothe/zuck-my-clothe-backend/docs"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
	"zuck-my-clothe/zuck-my-clothe-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title			Zuck-my-clothe API
// @version		1.0
// @description	This is API document for Zuck-my-clothe API
// @host			zuck-my-clothe-api.sokungz.work
//
//	@schemes		https
//
// @BasePath		/
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
	api.Get("/swaggerui/*", swagger.HandlerDefault)

	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8081",
		AllowCredentials: true,
	}))

	routes.RoutesRegister(db, api)

	fmt.Println(`
  ______          _    _____            
 |___  /         | |  |  __ \           
    / /_   _  ___| | _| |__) |_ _  __ _ 
   / /| | | |/ __| |/ /  ___/ _' |/ _' |
  / /_| |_| | (__|   <| |  | (_| | (_| |
 /_____\__,_|\___|_|\_\_|_  \__,_|\__,_|
  / ____|               (_)             
 | (___   ___ _ ____   ___  ___ ___     
  \___ \ / _ \ '__\ \ / / |/ __/ _ \    
  ____) |  __/ |   \ V /| | (_|  __/    
 |_____/ \___|_|    \_/ |_|\___\___|
	
 `)

	api.Listen(":" + config.PORT)
}
