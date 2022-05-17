package main

import (
	"log"

	"github.com/GorunovAlx/gophermart/internal/gophermart/app"
	"github.com/GorunovAlx/gophermart/internal/gophermart/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
