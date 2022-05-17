package app

import (
	"fmt"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
	"github.com/GorunovAlx/gophermart/pkg/logger"
	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())
	// Repository
	st := database.InitStorage(cfg)

	serviceShelf, err := v1.NewServiceShelf(cfg, st)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - serviceShelf.New: %w", err))
	}

	router := v1.Initialize(serviceShelf)

	router.Negroni.Run(cfg.RunAddress)
}
