package middleware

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// func SetupRedis(url string) *redis.Client {
func SetupRedis() *redis.Client {
	// Create a redis client.
	option, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to create a redis client: %v", err)
		return nil
	}
	client := redis.NewClient(option)

	ctx := context.Background()

	// Ping the Redis server and check if any errors occurred.
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	return client
}
