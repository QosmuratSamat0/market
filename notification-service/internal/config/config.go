package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `yaml:"env" env-default:"local" env:"ENV"`
	RabbitmqURL string `yaml:"rabbitmq_url" env:"RABBITMQ_URL" env-default:"amqp://guest:guest@rabbitmq:5672/"`
	HttpPort string `yaml:"http_port" env:"HTTP_PORT" env-default:"8085"`
}

func MustLoad() *Config {
	var cfg Config

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	configPath := fmt.Sprintf("./config/%s.yaml", env)

	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Printf("failed to read config: %v", err)
		}
	} else {
		cleanenv.ReadEnv(&cfg)
	}

	return &cfg
}
