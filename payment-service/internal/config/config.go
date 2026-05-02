package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string     `yaml:"env" env-default:"local" env:"ENV"`
	DatabaseURL     string     `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
	HTTPServer      HTTPServer `yaml:"http_server"`
	OrderServiceURL string     `yaml:"order_service_url" env:"ORDER_SERVICE_URL" env-default:"http://order-service:8083"`
	UserServiceURL  string     `yaml:"user_service_url" env:"USER_SERVICE_URL" env-default:"http://user-service:8081"`
	NatsURL         string     `yaml:"nats_url" env:"NATS_URL" env-default:"nats://nats:4222"`
	RabbitmqURL     string     `yaml:"rabbitmq_url" env:"RABBITMQ_URL" env-default:"amqp://guest:guest@rabbitmq:5672/"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8084" env:"HTTP_SERVER_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env-default:"15s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
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
			log.Printf("failed to read config from %s: %v", configPath, err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Printf("failed to read env: %v", err)
		}
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return &cfg
}
