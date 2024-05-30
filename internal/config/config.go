package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerHost         string        `envconfig:"SERVER_HOST" required:"true"`
	ServerPort         string        `envconfig:"SERVER_PORT" required:"true"`
	DockerRegistryHost string        `envconfig:"DOCKER_REGISTRY_HOST" required:"true"`
	DockerRegistryPort string        `envconfig:"DOCKER_REGISTRY_PORT" required:"true"`
	ClientTimeOut      time.Duration `envconfig:"CLIENT_TIMEOUT" required:"true"`
	CacheDefExp        time.Duration `envconfig:"CACHE_DEFAULT_EXPIRATION" required:"true"`
	CacheCleanInterval time.Duration `envconfig:"CACHE_CLEAN_UP_INTERVAL" required:"true"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("err to process config: %w", err)
	}

	return cfg, nil
}
