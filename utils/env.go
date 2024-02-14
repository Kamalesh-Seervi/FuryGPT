package utils

import (
	"github.com/kamalesh-seervi/simpleGPT/models"
	"github.com/spf13/viper"
)

func LoadConfig() (models.Config, error) {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		return models.Config{}, err
	}
	config := models.Config{
		APIKey:        viper.GetString("API_KEY"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		RedisDB:       viper.GetString("REDIS_DB"),
		RedisURL:      viper.GetString("REDIS_URL"),
	}

	return config, nil
}
