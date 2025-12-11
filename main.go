package main

import (
	"forger-companion/internal/app"
	"forger-companion/internal/config"
	"log"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		cfg = config.Default()
	}

	// Create and run app
	application := app.NewSimple(cfg)
	application.Run()
}
