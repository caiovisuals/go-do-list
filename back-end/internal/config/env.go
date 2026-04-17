package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL not set")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET not set")
	}
}
