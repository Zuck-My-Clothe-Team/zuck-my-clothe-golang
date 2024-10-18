package config

import (
	"os"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	FRONTEND_URL     string
	DB_DSN           string
	JWT_ACCESS_TOKEN string
	PORT             string
	APP_ENV          string
}

type RoutesRegister struct {
	DbConnection *platform.Postgres
	Config       *Config
	Application  *fiber.App
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	frontURL := os.Getenv("FRONTEND_URL")
	dbURL := os.Getenv("DB_DSN")
	jwtToken := os.Getenv("JWT_ACCESS_TOKEN")
	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")

	return &Config{
		FRONTEND_URL:     frontURL,
		DB_DSN:           dbURL,
		JWT_ACCESS_TOKEN: jwtToken,
		PORT:             port,
		APP_ENV:          appEnv,
	}, nil
}

func RouteRegister(db *platform.Postgres, config *Config, api *fiber.App) (*RoutesRegister, error) {

	if db == nil || config == nil || api == nil {
		panic("Error cannot create RouteRegister")
	}

	return &RoutesRegister{
		DbConnection: db,
		Config:       config,
		Application:  api,
	}, nil

}
