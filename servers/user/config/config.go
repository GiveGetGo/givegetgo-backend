package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init() {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("default")
	config.AddConfigPath("config/")
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("error on setting up default configuration: %v", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
