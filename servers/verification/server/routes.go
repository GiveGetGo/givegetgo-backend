package server

import (
	"verification/controller"
	"verification/middleware"
	"verification/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up verification utils
	verificationUtils := utils.NewVerificationUtils(DB, redisClient)
	defaultRateLimiter := middleware.SetupRateLimiter(redisClient, "60-M")
	sensitiveRateLimiter := middleware.SetupRateLimiter(redisClient, "10-M")

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	unAuthGroup.Use(defaultRateLimiter)
	{
		unAuthGroup.GET("/verification/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	verificationAuthGroup := r.Group("/v1/verification")
	verificationAuthGroup.Use(defaultRateLimiter)
	verificationAuthGroup.Use(middleware.AuthMiddleware())
	verificationAuthGroup.Use(sensitiveRateLimiter)
	{
		verificationAuthGroup.POST("/verify-email", controller.VerifyEmailVerificationHandler(verificationUtils))
	}

	// Internal routes - with auth middleware
	verificationInternalGroup := r.Group("/v1/internal/verification")
	verificationInternalGroup.Use(middleware.InternalAuthMiddleware())
	{
		verificationInternalGroup.POST("/request-email", controller.RequestEmailVerificationHandler(verificationUtils))
	}

	return r
}
