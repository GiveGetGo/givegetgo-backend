package main

import (
	"log"
	"os"
	"user_server/db"
	"user_server/middleware"
	"user_server/user"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.user")
	if err != nil {
		log.Fatalf("Error loading .env.user file: %v", err)
	}

	// load env
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

	// Internal routes
	internalRoutes := v1.Group("/internal")
	internalRoutes.Use(middleware.InternalAuthMiddleware())
	internalRoutes.POST("/user/email-verified", user.SetUserEmailVerifiedHandler(userUtils))

	// Public routes
	v1.POST("/user/register", user.UserRegisterHandler(userUtils))

	// start the server
	r.Run(":8080")
}
