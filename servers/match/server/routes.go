package server

import (
	"match/controller"
	"match/middleware"
	"match/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up match utils
	matchUtils := utils.NewMatchUtils(DB, redisClient)
	defaultRateLimiter := middleware.SetupRateLimiter(redisClient, "60-M")
	sensitiveRateLimiter := middleware.SetupRateLimiter(redisClient, "10-M")

	// Public routes - without auth middleware
	matchGroup := r.Group("/v1/match")
	matchGroup.Use(defaultRateLimiter)
	{
		matchGroup.GET("/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	matchAuthGroup := r.Group("/v1")
	matchAuthGroup.Use(defaultRateLimiter)
	matchAuthGroup.Use(middleware.AuthMiddleware())
	{
		defaultMatchAuthGroup := matchAuthGroup.Group("")
		{
			defaultMatchAuthGroup.GET("/match/:id", controller.GetMatchHandler(matchUtils))
		}

		sensitiveMatchAuthGroup := matchAuthGroup.Group("")
		sensitiveMatchAuthGroup.Use(sensitiveRateLimiter)
		{
			sensitiveMatchAuthGroup.POST("/match", controller.MatchHandler(matchUtils))
		}
	}

	return r
}
