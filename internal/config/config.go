package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}