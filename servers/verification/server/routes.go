package server

import (
	"verification/controller"
	"verification/middleware"
	"verification/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up verification utils
	verificationUtils := utils.NewVerificationUtils(DB, redisClient)

	// Public routes - without auth middleware

	// Public routes - with auth middleware
	verificationAuthGroup := r.Group("/v1/verification")
	// TODO: Add auth middleware
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
