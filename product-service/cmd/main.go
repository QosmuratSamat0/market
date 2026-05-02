package main

import (
	"log"

	"github.com/QosmuratSamat0/product-service/internal/app"
	"github.com/QosmuratSamat0/product-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	application.Run()
}
