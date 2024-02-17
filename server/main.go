package main

import (
	"log"
	"os"
	"server/db"
	"server/middleware"
	"server/user"
	"server/utils"

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
	DB := db.ConnectPostgresDB(dburl)
	db.AutoMigratePostgresDB(DB)
	userUtils := utils.NewUserUtils(DB)
	rateLimiter := middleware.SetupRateLimiter()

	// set up the router and v1 routes
	r := gin.Default()
	v1 := r.Group("/v1")
	v1.Use(rateLimiter)

	// routes
	v1.POST("/user/register", user.UserRegisterHandler(userUtils))

	// start the server
	r.Run(":8080")
}
