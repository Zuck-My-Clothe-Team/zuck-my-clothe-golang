package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FRONTEND_URL     string
	DB_DSN           string
	JWT_ACCESS_TOKEN string
	PORT             string
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

	return &Config{
		FRONTEND_URL:     frontURL,
		DB_DSN:           dbURL,
		JWT_ACCESS_TOKEN: jwtToken,
		PORT:             port,
	}, nil
}
