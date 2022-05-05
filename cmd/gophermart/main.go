package main

import (
	"log"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/app"
	"github.com/GorunovAlx/gophermart/pkg/migration"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Migration
	migration.Run(cfg)

	// Run
	app.Run(cfg)
}
