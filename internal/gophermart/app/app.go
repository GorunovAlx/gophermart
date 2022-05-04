package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/GorunovAlx/gophermart/pkg/httpserver"
	"github.com/GorunovAlx/gophermart/pkg/logger"
	"github.com/GorunovAlx/gophermart/pkg/postgres"
	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())

	// Repository
	pg, err := postgres.New(cfg, postgres.MaxPoolSize(2))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	serviceShelf, err := v1.NewServiceShelf(cfg, pg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	router := v1.NewHandler(serviceShelf)
	router.InitializeRoutes()

	httpServer := httpserver.New(router.Negroni, httpserver.Port(cfg.RunAddress))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
