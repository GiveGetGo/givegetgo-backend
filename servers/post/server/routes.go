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

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	{
		unAuthGroup.GET("/post/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	postAuthGroup := r.Group("/v1")
	postAuthGroup.Use(middleware.AuthMiddleware())
	{
		postAuthGroup.POST("/post", controller.AddPostHandler(postUtils))
		postAuthGroup.GET("/post", controller.GetPostHandler(postUtils))
		postAuthGroup.GET("/post/archive", controller.GetPostArchiveHandler(postUtils))
		postAuthGroup.GET("/post/by-user", controller.GetPostByUserIdHandler(postUtils))
		postAuthGroup.GET("/post/:id", controller.GetPostByPostIdHandler(postUtils))
		postAuthGroup.PUT("/post/:id", controller.EditPostByIdHandler(postUtils))
		postAuthGroup.DELETE("/post/:id", controller.DeletePostHandler(postUtils))
	}

	// interal routes
	postInternalGroup := r.Group("/v1/internal")
	postInternalGroup.Use(middleware.InternalAuthMiddleware())
	{
		postInternalGroup.PUT("/post/status", controller.UpdatePostStatusHandler(postUtils))
	}

	return r
}
