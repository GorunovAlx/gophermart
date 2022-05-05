package app

import (
	"fmt"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/pkg/logger"
	"github.com/GorunovAlx/gophermart/pkg/postgres"
	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())

	l.Debug("cfg: runaddress: %v, database: %v, accrual: %v", cfg.RunAddress, cfg.DatabaseURI, cfg.AccrualAddress)
	// Repository
	pg, err := postgres.New(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	l.Debug("postgres: %v", pg.Pool)
	defer pg.Close()

	//serviceShelf, err := v1.NewServiceShelf(cfg, pg)
	serviceShelf, err := v1.NewServiceShelf(cfg, nil)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - serviceShelf.New: %w", err))
	}

	router := v1.Initialize(serviceShelf)

	router.Negroni.Run(cfg.RunAddress)
}
