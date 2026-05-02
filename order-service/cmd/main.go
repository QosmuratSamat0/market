package main

import (
	"log"

	"github.com/QosmuratSamat/order-service/internal/app"
	"github.com/QosmuratSamat/order-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	application.Run()
}
