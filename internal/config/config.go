package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	JWTSecret             string `env:"JWT_SECRET" required:"true"`
	ApiServiceServerPort  string `env:"API_SERVICE_SERVER_PORT" required:"true"`
	UserServiceAddress    string `env:"USER_SERVICE_ADDRESS" required:"true"`
	AuthServiceAddress    string `env:"AUTH_SERVICE_ADDRESS" required:"true"`
	PostServiceAddress    string `env:"POST_SERVICE_ADDRESS" required:"true"`
	StorageServiceAddress string `env:"STORAGE_SERVICE_ADDRESS" required:"true"`
}

func Load() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
