package main

import (
	"log"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/app"
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
