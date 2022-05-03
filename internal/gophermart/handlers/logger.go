package handlers

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/GorunovAlx/gophermart/internal/gophermart/config"
)

var Logger zerolog.Logger

func CreateLogger() {
	logfile, err := os.OpenFile("server_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot start %v", logfile)
	}
	//defer logfile.Close()

	Logger = zerolog.New(logfile).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.Level(config.Cfg.ZerologLevel))

	Logger.Info().Msg("start server")
}
