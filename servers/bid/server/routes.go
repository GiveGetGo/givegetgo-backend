package server

import (
	"bid/controller"
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

	// Public routes - without auth middleware
	bidGroup := r.Group("/v1/bid")
	{
		bidGroup.GET("/health", sharedController.HealthCheckHandler())
		bidGroup.POST("/bid", controller.AddBidHandler(bidUtils))
		bidGroup.GET("/getbidbyp/:postID", controller.GetBidsForPostHandler(bidUtils))
		bidGroup.GET("/getbidbyb/:bidID", controller.FindBidByIDHandler(bidUtils))
		bidGroup.DELETE("/delete/:bidID", controller.DeleteBidHandler(bidUtils))
		bidGroup.PUT("/put/:bidID", controller.UpdateBidDescriptionHandler(bidUtils))

	}

	return r
}
