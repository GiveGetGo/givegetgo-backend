package main

import (
	"notification/server"

	"github.com/GiveGetGo/shared/config"
)

func main() {
	config.LoadEnv(".env.notification")

	server.StartServer()
}
