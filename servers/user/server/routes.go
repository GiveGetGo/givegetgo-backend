package server

import (
	"user/controller"
	"user/middleware"
	"user/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up user utils
	userUtils := utils.NewUserUtils(DB, redisClient)

	// Set up rate limiter
	rateLimiter := middleware.SetupRateLimiter(redisClient)

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	unAuthGroup.Use(rateLimiter) // Apply rate limiter
	{
		unAuthGroup.POST("/user/register", controller.RegisterHandler(userUtils))
		unAuthGroup.POST("/user/login", controller.LoginHandler(userUtils))
	}

	// Public routes - with auth middleware
	authGroup := r.Group("/v1")
	authGroup.Use(rateLimiter) // Apply rate limiter
	{
		userGroup := authGroup.Group("/user")
		{
			userGroup.POST("/forgot-password", controller.ForgotPasswordHandler(userUtils))
		}

		mfaGroup := authGroup.Group("/mfa")
		{
			mfaGroup.POST("/request", controller.RequestMFAVerificationHandler(userUtils))
			mfaGroup.POST("/verify", controller.VerifyMFAHandler(userUtils))
			mfaGroup.GET("/qrcode/:userid", controller.MFAQRCodeHandler(userUtils))
		}
	}

	// Internal routes - with auth middleware
	internalGroup := r.Group("/v1/internal")
	{
		internalGroup.POST("/user/email-verified", controller.SetUserEmailVerifiedHandler(userUtils))
	}

	return r
}
