package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load environment variables
	err := godotenv.Load(".env.verification")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
