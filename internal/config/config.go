package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	RedisPort string `env:"REDIS_PORT" required:"true"`

	JWTSecret            string `env:"JWT_SECRET" required:"true"`
	ApiServiceServerPort string `env:"API_SERVICE_SERVER_PORT" required:"true"`
}

func Load() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
