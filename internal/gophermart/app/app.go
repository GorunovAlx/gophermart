package app

import (
	"fmt"
	"net/http"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/pkg/logger"
	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())

	l.Debug("cfg: runaddress: %v, database: %v, accrual: %v", cfg.RunAddress, cfg.DatabaseURI, cfg.AccrualAddress)
	// Repository
	/*
		pg, err := postgres.New(cfg)
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
		}
		l.Debug("postgres: %v", pg.Pool)
		defer pg.Close()
	*/
	serviceShelf, err := v1.NewServiceShelf(cfg, nil)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - serviceShelf.New: %w", err))
	}

	router := v1.Initialize(serviceShelf)

	l.Fatal(http.ListenAndServe(cfg.RunAddress, router.Negroni))

	//router.Negroni.Run(cfg.RunAddress)

	/*
		s := &http.Server{
			Addr:         cfg.RunAddress,
			Handler:      router.Negroni,
			WriteTimeout: 1015 * time.Second,
			ReadTimeout:  1015 * time.Second,
			IdleTimeout:  time.Second * 60 * 5,
		}
		l.Fatal(s.ListenAndServe())
	*/
}
