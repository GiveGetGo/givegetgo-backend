/*package server

import (
	"user/db"
	"user/middleware"
)

func StartServer() {
	DB := db.InitDB()                      // Initialize the database
	redisClient := middleware.SetupRedis() // Set up Redis

	r := NewRouter(DB, redisClient) // Set up the router and v1 routes
	r.Run(":8080")                  // Start the server
}
*/

package server

import (
    "time"
    "user/db"
    "user/middleware"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func StartServer() {
    DB := db.InitDB()                      // Initialize the database
    redisClient := middleware.SetupRedis() // Set up Redis

    r := gin.Default()

    // Configure CORS middleware options as needed
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // "http://128.210.106.61:8080" or "*" for testing purposes
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    r = NewRouter(DB, redisClient) // Set up the router and v1 routes (assuming NewRouter sets up the routes)
    r.Run(":8080")                 // Start the server
}