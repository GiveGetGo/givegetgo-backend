package main

import (
	"log"
	"os"
	"server/db"

	"github.com/joho/godotenv"
)

func main() {
	// load the .env file
	err := godotenv.Load(".env.server")
	if err != nil {
		log.Fatalf("Error loading .env.server file: %v", err)
	}

	// load db variables
	dburl := os.Getenv("DATABASE_URL")

	// Set up database connection
	_, err = db.ConnectPostgresDB(dburl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//
}
