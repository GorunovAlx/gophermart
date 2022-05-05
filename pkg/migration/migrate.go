package migration

import (
	"errors"
	"log"
	"time"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func Run(cfg *config.Config) {
	if len(cfg.DatabaseURI) == 0 {
		log.Fatalf("migrate: environment variable not declared: DATABASE_URI")
	}

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", cfg.DatabaseURI)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migrate: up error: %s", err)
		}

		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrate: no change")
			return
		}
	}

	log.Printf("Migrate: up success")
}
