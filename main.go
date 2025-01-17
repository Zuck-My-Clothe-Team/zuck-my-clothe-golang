package main

import (
	"fmt"
	"log"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	nacronsritammarat "zuck-my-clothe/zuck-my-clothe-backend/cron"
	"zuck-my-clothe/zuck-my-clothe-backend/docs"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"
	"zuck-my-clothe/zuck-my-clothe-backend/routes"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title			Zuck-my-clothe API
// @version		1.0
// @description	This is API document for Zuck-my-clothe API
// @BasePath		/
func main() {
	time.Local = time.UTC
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Can't load config", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	db, dbErr := platform.InitDB(cfg.DB_DSN)

	if dbErr != nil {
		log.Fatal("Can not Init Database", dbErr)
	}

	valErr := validatorboi.CreateValidator()

	if valErr != "success" {
		log.Fatal("Can not Init Validator")
	}

	konCron := nacronsritammarat.SummonKonCron(db)
	konCron.StartKonKron()

	api := fiber.New()

	if cfg.APP_ENV == "PRODUCTION" {
		docs.SwaggerInfo.Host = "zuck-my-clothe-api.sokungz.work"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:3000"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:3001, http://localhost:8081, https://zuck-my-clothe.sokungz.work",
		AllowCredentials: true,
	}))

	routeRegister, err := config.RouteRegister(db, cfg, api)

	if err != nil {
		log.Fatal("Cannot initial route register", err)
	}
	routes.RoutesRegister(routeRegister)

	api.Get("/swaggerui/*", swagger.HandlerDefault)

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

	api.Listen(":" + cfg.PORT)
}
