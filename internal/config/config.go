package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	DriverServiceClient struct {
		URL    string `mapstructure:"url"`
		ApiKey string `mapstructure:"api_key"`
	} `mapstructure:"driver_service_client"`
	JWTSecretKey string `mapstructure:"jwt_secret_key"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
