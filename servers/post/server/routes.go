package server

import (
	"post/controller"
	"post/middleware"
	"post/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up verification utils
	postUtils := utils.NewPostUtils(DB, redisClient)
	defaultRateLimiter := middleware.SetupRateLimiter(redisClient, "60-M")
	sensitiveRateLimiter := middleware.SetupRateLimiter(redisClient, "10-M")

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	unAuthGroup.Use(defaultRateLimiter)
	{
		unAuthGroup.GET("/post/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	postAuthGroup := r.Group("/v1")
	postAuthGroup.Use(defaultRateLimiter)
	postAuthGroup.Use(middleware.AuthMiddleware())
	{
		defaultPostAuthGroup := postAuthGroup.Group("")
		{
			defaultPostAuthGroup.GET("/post", controller.GetPostHandler(postUtils))
			defaultPostAuthGroup.GET("/post/:id", controller.GetPostByPostIdHandler(postUtils))
			defaultPostAuthGroup.GET("/post/by-user", controller.GetPostByUserIdHandler(postUtils))
		}

		sensitivePostAuthGroup := postAuthGroup.Group("")
		sensitivePostAuthGroup.Use(sensitiveRateLimiter)
		{
			sensitivePostAuthGroup.POST("/post", controller.AddPostHandler(postUtils))
			sensitivePostAuthGroup.GET("/post/archive", controller.GetPostArchiveHandler(postUtils))
			sensitivePostAuthGroup.PUT("/post/:id", controller.EditPostByIdHandler(postUtils))
			sensitivePostAuthGroup.DELETE("/post/:id", controller.DeletePostHandler(postUtils))
		}
	}

	// interal routes
	postInternalGroup := r.Group("/v1/internal")
	postInternalGroup.Use(middleware.InternalAuthMiddleware())
	{
		postInternalGroup.PUT("/post/status", controller.UpdatePostStatusHandler(postUtils))
	}

	return r
}
