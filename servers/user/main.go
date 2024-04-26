package main

import (
	"user/config"
	"user/server"

	sharedConfig "github.com/GiveGetGo/shared/config"
)

func main() {
	sharedConfig.LoadEnv(".env.user") // Load environment variables

	config.Init()        // Initialize Config
	server.StartServer() // Start the server
}
