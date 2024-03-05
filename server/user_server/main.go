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

	// Set up Redis
	redisClient := middleware.SetupRedis()

	// Set up rate limiter
	rateLimiter := middleware.SetupRateLimiter(redisClient)

	// Set up user utils
	userUtils := utils.NewUserUtils(DB, redisClient)

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
	v1.POST("/user/login", user.LoginHandler(userUtils))
	v1.POST("/user/request-mfa", user.RequestMFAVerificationHandler(userUtils))
	v1.POST("/user/verify-mfa", user.VerifyMFAHandler(userUtils))
	v1.GET("/user/mfa-qrcode/:userid", user.MFAQRCodeHandler(userUtils))
	v1.POST("user/forgot-password", user.ForgotPasswordHandler(userUtils))

	// start the server
	r.Run(":8080")
}
