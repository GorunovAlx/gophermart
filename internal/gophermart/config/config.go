package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	// The server address
	RunAddress string `env:"RUN_ADDRESS"`
	// Database url string like postgres://postgres:pass@localhost:5432/dbname?sslmode=disable
	DatabaseURI string `env:"DATABASE_URI"`
	// The accrual address
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	// logging level for zerolog
	ZerologLevel int8 `env:"ZERO_LOG_LEVEL" envDefault:"0"`
	// logging level for pgx driver db
	PgxLogLevel string `env:"PGX_LOG_LEVEL" envDefualt:"info"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	parameters := os.Args[1:]

	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}

	if len(parameters) > 0 {
		if flag.Lookup("a") == nil {
			flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "HTTP server launch address")
		}
		if flag.Lookup("d") == nil {
			flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "Database url string")
		}
		if flag.Lookup("r") == nil {
			flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "Addres of accrual service")
		}
		flag.Parse()
	}

	return cfg, nil
}
