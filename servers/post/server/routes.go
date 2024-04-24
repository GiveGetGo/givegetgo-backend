package server

import (
	"post/controller"
	"post/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up verification utils
	postUtils := utils.NewPostUtils(DB, redisClient)

	// Public routes - without auth middleware

	// Public routes - with auth middleware
	postAuthGroup := r.Group("/v1/post")
	// TODO: Add auth middleware
	{
		postAuthGroup.POST("/add", controller.AddPostHandler(postUtils))
		postAuthGroup.GET("/getp", controller.GetPostHandler(postUtils))
		postAuthGroup.DELETE("/delete", controller.DeletePostHandler(postUtils))
	}

	return r
}