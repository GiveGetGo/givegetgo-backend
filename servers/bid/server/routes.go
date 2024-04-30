package server

import (
	"bid/controller"
	"bid/middleware"
	"bid/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up match utils
	bidUtils := utils.NewBidUtils(DB, redisClient)
	defaultRateLimiter := middleware.SetupRateLimiter(redisClient, "60-M")
	sensitiveRateLimiter := middleware.SetupRateLimiter(redisClient, "10-M")

	// Public routes - without auth middleware
	bidGroup := r.Group("/v1/bid")
	bidGroup.Use(defaultRateLimiter)
	{
		bidGroup.GET("/health", sharedController.HealthCheckHandler())
	}

	bidAuthGroup := r.Group("/v1")
	bidAuthGroup.Use(defaultRateLimiter)
	bidAuthGroup.Use(middleware.AuthMiddleware())
	{
		defaultBidAuthGroup := bidAuthGroup.Group("")
		{
	        defaultBidAuthGroup.GET("/by-post/:postid", controller.GetBidsForPostHandler(bidUtils))
		    defaultBidAuthGroup.GET("/:bidid", controller.FindBidByIDHandler(bidUtils))
		}
		sensitiveBidAuthGroup := bidAuthGroup.Group("")
		sensitiveBidAuthGroup.Use(sensitiveRateLimiter)
		{
			sensitiveBidAuthGroup.POST("/bid", controller.AddBidHandler(bidUtils))
            sensitiveBidAuthGroup.POST("/by-post/:postid", controller.AddBidHandler(bidUtils))
		    sensitiveBidAuthGroup.DELETE("/:bidid", controller.DeleteBidHandler(bidUtils))
		    sensitiveBidAuthGroup.PUT("/:bidid", controller.UpdateBidDescriptionHandler(bidUtils))
		}
	}

	return r
}
