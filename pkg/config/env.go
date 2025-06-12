package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabaseURL   string
	APIGatewayKey string
	Port          string
	LogLevel      string
	AppEnv        string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		APIGatewayKey: os.Getenv("API_GATEWAY_KEY"),
		Port:          GetEnv("PORT", "8080"),
		LogLevel:      GetEnv("LOG_LEVEL", "info"),
		AppEnv:        GetEnv("APP_ENV", "development"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.APIGatewayKey == "" {
		log.Warn().Msg("API_GATEWAY_KEY not set - API authentication disabled")
	}

	return nil
}
