package main

import (
	"user/server"

	"github.com/GiveGetGo/shared/config"
)

func main() {
	config.LoadEnv(".env.user") // Load environment variables

	server.StartServer() // Start the server
}


/*
package main

import (
	"time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Configure CORS middleware options as needed
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://128.210.106.61:8080"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // your routes here

    r.Run(":8080") // Change this to your backend's actual port
}
*/