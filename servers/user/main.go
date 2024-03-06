package main

import (
	"user/server"

	"github.com/GiveGetGo/shared/config"
)

func main() {
	config.LoadEnv(".env.user") // Load environment variables

	server.StartServer() // Start the server
}
