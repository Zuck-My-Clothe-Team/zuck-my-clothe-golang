package main

import (
	"fmt"
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

	api := fiber.New()
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8081",
		AllowCredentials: true,
	}))

	api.Listen(":3000")
}
