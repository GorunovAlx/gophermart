package app

import (
	"github.com/GorunovAlx/gophermart/internal/gophermart/app/accrual"
	"github.com/GorunovAlx/gophermart/internal/gophermart/app/logger"
	"github.com/GorunovAlx/gophermart/internal/gophermart/config"
	"github.com/GorunovAlx/gophermart/internal/gophermart/database"
	"github.com/GorunovAlx/gophermart/internal/gophermart/entity"

	"github.com/rs/zerolog"

	v1 "github.com/GorunovAlx/gophermart/internal/gophermart/controller/http/v1"
)

func Run(cfg *config.Config) {
	l := logger.New(zerolog.Level(cfg.ZerologLevel).String())
	// Repository
	st := database.InitStorage(cfg)
	us := entity.UserStorage{S: *st}
	os := entity.OrderStorage{S: *st}
	ws := entity.WithdrawStorage{S: *st}
	as := accrual.NewAccrualService(cfg.AccrualAddress, os)

	router := v1.NewHandler(us, os, ws, as, l)

	router.Negroni.Run(cfg.RunAddress)
}
