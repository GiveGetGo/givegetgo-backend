package main

import (
	"log"
	"os"
	"verification_server/db"
	"verification_server/middleware"
	"verification_server/utils"
	"verification_server/verification"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.verification")
	if err != nil {
		log.Fatalf("Error loading .env.verification file: %v", err)
	}

	// load db variables
	dburl := os.Getenv("DATABASE_URL")

	// Set up database connection
	DB := db.ConnectPostgresDB(dburl)
	db.AutoMigratePostgresDB(DB)
	verificationUtils := utils.NewVerificationUtils(DB)

	// set up the router and v1 routes
	r := gin.Default()
	v1 := r.Group("/v1")
	v1.Use(middleware.InternalAuthMiddleware())

	// routes
	v1.POST("/verification/request-email", verification.RequestEmailVerificationHandler(verificationUtils))
	v1.POST("/verification/verify-email", verification.VerifyEmailVerificationHandler(verificationUtils))

	// start the server
	r.Run(":8080")
}
