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
	DatabaseURI string `env:"DATABASE_URI" envDefault:""`
	//
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:":8080"`
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
			flag.StringVar(&Cfg.DatabaseURI, "f", Cfg.DatabaseURI, "")
		}
		if flag.Lookup("r") == nil {
			flag.StringVar(&Cfg.AccrualAddress, "d", Cfg.AccrualAddress, "")
		}
		flag.Parse()
	}
}
