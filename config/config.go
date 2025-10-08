package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DATABASE_URL string
	REDIS_URL string
	AIRPORT_API_URL string
	WEATHER_API_URL string
	WEATHER_API_KEY string
}

func Load() (Config, error) {
	var config Config
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	return config, err
}