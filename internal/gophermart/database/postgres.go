package database

import (
	"context"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	PGpool *pgx.Conn
}

func InitStorage(cfg *config.Config) *Storage {
	makeMigration(cfg.DatabaseURI)
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}

	return &Storage{
		PGpool: conn,
	}
}

func makeMigration(uri string) {
	m, err := migrate.New(
		"file://internal/gophermart/database/migrations",
		uri)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
