package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// SetupRateLimiter configures rate limiting with a specified rate.
func SetupRateLimiter(client *redis.Client, rateLimit string) gin.HandlerFunc {
	// Define a limit rate from the rateLimit parameter.
	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		log.Fatalf("Failed to create a rate: %v", err)
		return nil
	}

	// Create a store with the Redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "givegetgo-limiter",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create a store: %v", err)
		return nil
	}

	// Create a new middleware with the limiter instance.
	middleware := mgin.NewMiddleware(limiter.New(store, rate))

	return middleware
}
