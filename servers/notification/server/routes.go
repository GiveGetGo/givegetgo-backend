package server

import (
	"notification/controller"
	"notification/middleware"
	"notification/utils"

	sharedController "github.com/GiveGetGo/shared/controller"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(DB *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()

	// Set up verification utils
	notificationUtils := utils.NewNotificationUtils(DB, redisClient)

	// Public routes - without auth middleware
	unAuthGroup := r.Group("/v1")
	{
		unAuthGroup.GET("/notification/health", sharedController.HealthCheckHandler())
	}

	// Public routes - with auth middleware
	notificationAuthGroup := r.Group("/v1")
	notificationAuthGroup.Use(middleware.AuthMiddleware())
	{
		notificationAuthGroup.GET("/notfication/:id", controller.GetNotificationByUserID(notificationUtils))
		notificationAuthGroup.DELETE("/notification/:id", controller.DeleteNotification(notificationUtils))
	}

	// interal routes
	notificationInternalGroup := r.Group("/v1/internal")
	notificationInternalGroup.Use(middleware.InternalAuthMiddleware())
	{
		notificationInternalGroup.POST("/notification", controller.CreateNewNotification(notificationUtils))
	}

	return r
}
