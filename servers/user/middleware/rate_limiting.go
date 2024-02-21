package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func SetupRateLimiter(client *redis.Client) gin.HandlerFunc {
	// Define a limit rate to 4 request per minute.
	rate, err := limiter.NewRateFromFormatted("8-M")
	if err != nil {
		log.Fatalf("Failed to create a rate: %v", err)
		return nil
	}

	// Create a store with the redis client.
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

// TODO: setup different rate limiters for different routes
