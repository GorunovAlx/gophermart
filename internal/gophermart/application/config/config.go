package config

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	// The server address
	RunAddress string `env:"RUN_ADDRESS" envDefault:":8081"`
	// Database url string like postgres://postgres:pass@localhost:5432/dbname?sslmode=disable
	// postgres://postgres:barkleys@localhost:5432/dbgophermart?sslmode=disable
	DatabaseURI string `env:"DATABASE_URI"`
	//
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:":8080"`
	// logging level for zerolog
	ZerologLevel int8 `env:"ZERO_LOG_LEVEL" envDefault:"0"`
	// logging level for pgx driver db
	PgxLogLevel string `env:"PGX_LOG_LEVEL" envDefualt:"info"`
}

var Cfg Config

func SetConfig() {
	parameters := os.Args[1:]

	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}

	if len(parameters) > 0 {
		if flag.Lookup("a") == nil {
			flag.StringVar(&Cfg.RunAddress, "a", Cfg.RunAddress, "HTTP server launch address")
		}
		if flag.Lookup("d") == nil {
			flag.StringVar(&Cfg.DatabaseURI, "f", Cfg.DatabaseURI, "postgres://postgres:pass@localhost:5432/dbname?sslmode=disable")
		}
		if flag.Lookup("r") == nil {
			flag.StringVar(&Cfg.AccrualAddress, "d", Cfg.AccrualAddress, "Addres of accrual service")
		}
		flag.Parse()
	}
}
