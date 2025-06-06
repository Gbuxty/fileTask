package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTP HTTPServer
}

func NewConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{HTTP: newHTTPServer()}, nil
}
