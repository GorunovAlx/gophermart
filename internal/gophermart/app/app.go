package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/pkg/logger"
	"github.com/GorunovAlx/gophermart/pkg/postgres"
	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())

	// Repository
	pg, err := postgres.New(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	serviceShelf, err := v1.NewServiceShelf(cfg, pg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	router := v1.NewHandler(serviceShelf)

	s := &http.Server{
		Addr:           cfg.RunAddress,
		Handler:        router.Negroni,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	l.Fatal(s.ListenAndServe())
}
