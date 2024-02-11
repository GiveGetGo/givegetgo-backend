package main

import (
	"log"
	"os"
	"server/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	// set up the router and v1 routes
	r := gin.Default()
	v1 := r.Group("/v1")

	// start the server
	r.Run(":8080")
}
