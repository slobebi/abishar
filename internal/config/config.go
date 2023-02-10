package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTPServer HTTPServer
		Database   Database
		Redis      Redis
		JWT        JWT
	}

	HTTPServer struct {
		ListenAddress   string
		Port            int
		GracefulTimeout time.Duration
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
	}

	Database struct {
		User     string
		Password string
		DBName   string
		Host     string
		Port     string
		SSLMode  string
		Retry    int
	}

	Redis struct {
		Endpoint string
		Timeout  int
		MaxIdle  int
	}

	JWT struct {
		Secret string
	}
)

func Init() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("../../")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
