package main

import (
	"github.com/spf13/viper"
)

func main() {
	loadConfig()
}

func loadConfig() {
	viper.SetConfigFile("config.yml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
