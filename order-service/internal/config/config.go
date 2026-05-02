package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                  string     `yaml:"env" env-default:"local" env:"ENV"`
	DatabaseURL          string     `yaml:"conn_db" env:"DATABASE_URL"`
	HTTPServer           HTTPServer `yaml:"http_server"`
	JWTSecret            string     `yaml:"jwt_secret" env:"JWT_SECRET"`
	ProductServiceURL    string     `yaml:"product_service_url" env:"PRODUCT_SERVICE_URL" env-default:"http://product-service:8082"`
	NatsURL              string     `yaml:"nats_url" env:"NATS_URL" env-default:"nats://nats:4222"`
	AppURL               string     `env:"APP_URL"`
	AllowLocalhostOrigin string     `yaml:"ALLOW_LOCALHOST_ORIGIN" env:"ALLOW_LOCALHOST_ORIGIN"`

	AllowedOrigins []string
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8083" env:"HTTP_SERVER_ADDRESS"`
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
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	var origins []string
	if cfg.AppURL != "" {
		origins = append(origins, cfg.AppURL)
	}
	if cfg.AllowLocalhostOrigin != "" {
		for _, o := range strings.Split(cfg.AllowLocalhostOrigin, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
	}
	cfg.AllowedOrigins = origins

	return &cfg
}
