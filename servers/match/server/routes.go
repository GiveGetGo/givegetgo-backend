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

	// Public routes - without auth middleware
	matchGroup := r.Group("/v1/match")
	{
		matchGroup.GET("/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	matchAuthGroup := r.Group("/v1")
	matchAuthGroup.Use(middleware.AuthMiddleware())
	{
		matchGroup.POST("/match", controller.MatchHandler(matchUtils))
		matchGroup.GET("/match/:id", controller.GetMatchHandler(matchUtils))
	}

	return r
}
