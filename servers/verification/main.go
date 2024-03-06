package main

import (
	"verification/server"

	"github.com/GiveGetGo/shared/config"
)

func main() {
	config.LoadEnv(".env.verification") // Load environment variables

	server.StartServer() // Start the server
}
