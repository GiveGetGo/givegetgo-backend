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

	userUtils := utils.NewUserUtils(DB, redisClient) // Set up user utils
	defaultRateLimiter := middleware.SetupRateLimiter(redisClient, "60-M")
	sensitiveRateLimiter := middleware.SetupRateLimiter(redisClient, "10-M")

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	unAuthGroup.Use(defaultRateLimiter)
	{
		userUnAuthGroup := unAuthGroup.Group("")
		{
			userUnAuthGroup.GET("/user/health", sharedController.HealthCheckHandler())
			userUnAuthGroup.GET("/user/logout", controller.LogoutHandler(userUtils))
		}

		sensitiveUnAuthGroup := unAuthGroup.Group("")
		sensitiveUnAuthGroup.Use(sensitiveRateLimiter)
		{
			sensitiveUnAuthGroup.POST("/user/register", controller.RegisterHandler(userUtils))
			sensitiveUnAuthGroup.POST("/user/login", controller.LoginHandler(userUtils))
		}
	}

	// Public routes - with auth middleware
	authGroup := r.Group("/v1")
	authGroup.Use(defaultRateLimiter)
	authGroup.Use(middleware.AuthMiddleware())
	{
		userGroup := authGroup.Group("/user")
		{
			userGroup.GET("/session", controller.SessionHandler(userUtils))
			userGroup.GET("/verified", controller.VerifiedHandler(userUtils))
			userGroup.GET("/me", controller.GetMeHandler(userUtils))
			userGroup.PUT("/me", controller.EditMeHandler(userUtils))
		}

		sensitiveUserGroup := userGroup.Group("")
		sensitiveUserGroup.Use(sensitiveRateLimiter)
		{
			sensitiveUserGroup.POST("/forgot-password", controller.ForgotPasswordHandler(userUtils))
			sensitiveUserGroup.POST("/reset-password", controller.ResetPasswordHandler(userUtils))
		}

		mfaGroup := authGroup.Group("/mfa")
		mfaGroup.Use(sensitiveRateLimiter)
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
