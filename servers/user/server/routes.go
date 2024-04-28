package server

import (
	"user/config"
	"user/controller"
	"user/middleware"
	"user/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	store := config.InitSession()                // Initialize session store using config
	r.Use(sessions.Sessions("givegetgo", store)) // Use sessions with the store

	userUtils := utils.NewUserUtils(DB, redisClient)        // Set up user utils
	rateLimiter := middleware.SetupRateLimiter(redisClient) // Setup Rate limiter

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	unAuthGroup.Use(rateLimiter) // Apply rate limiter
	{
		unAuthGroup.GET("/user/health", sharedController.HealthCheckHandler())
		unAuthGroup.POST("/user/register", controller.RegisterHandler(userUtils))
		unAuthGroup.POST("/user/login", controller.LoginHandler(userUtils))
	}

	// Public routes - with auth middleware
	authGroup := r.Group("/v1")
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.Use(rateLimiter) // Apply rate limiter
	{
		userGroup := authGroup.Group("/user")
		{
			userGroup.GET("/session", controller.SessionHandler(userUtils))
			userGroup.GET("/verified", controller.VerifiedHandler(userUtils))
			userGroup.GET("/me", controller.GetMeHandler(userUtils))
			userGroup.PUT("/me", controller.EditMeHandler(userUtils))
			userGroup.POST("/forgot-password", controller.ForgotPasswordHandler(userUtils))
			userGroup.POST("/reset-password", controller.ResetPasswordHandler(userUtils))
		}

		mfaGroup := authGroup.Group("/mfa")
		{
			mfaGroup.POST("", controller.VerifyMFAHandler(userUtils))
			mfaGroup.GET("", controller.GetMFAHandler(userUtils))
		}
	}

	// Internal routes - with auth middleware
	internalGroup := r.Group("/v1/internal")
	internalGroup.Use(middleware.InternalAuthMiddleware())
	{
		internalGroup.POST("/user/email-verified", controller.SetUserEmailVerifiedHandler(userUtils))
	}

	return r
}
