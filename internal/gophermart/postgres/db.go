package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GorunovAlx/gophermart/internal/gophermart/application/config"

	//"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreatePGXPool(ctx context.Context, dsn string, logger pgx.Logger, logLevel pgx.LogLevel) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conf.ConnConfig.Logger = logger

	if logLevel != 0 {
		conf.ConnConfig.LogLevel = logLevel
	}

	conf.MaxConns = 20
	conf.MaxConnIdleTime = time.Second * 30
	conf.MaxConnLifetime = time.Minute * 2

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, nil
}

// LogLevelFromEnv returns the pgx.LogLevel from the environment variable PGX_LOG_LEVEL.
// By default this is info (pgx.LogLevelInfo), which is good for development.
func LogLevelFromEnv() (pgx.LogLevel, error) {
	if level := config.Cfg.PgxLogLevel; level != "" {
		l, err := pgx.LogLevelFromString(level)
		if err != nil {
			return pgx.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return pgx.LogLevelInfo, nil
}

// PGXStdLogger prints pgx logs to the standard logger.
// os.Stderr by default.
type PGXStdLogger struct{}

func (l *PGXStdLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

func NewDBStorage() (*pgxpool.Pool, error) {
	pgxLogLevel, err := LogLevelFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	/*
		m, err := migrate.New(
			"file:///migrations",
			config.Cfg.DatabaseURI,
		)
		if err != nil {
			panic(err)
		}
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
	*/

	pgPool, err := CreatePGXPool(context.Background(), config.Cfg.DatabaseURI, &PGXStdLogger{}, pgxLogLevel)
	if err != nil {
		log.Fatal(err)
	}

	conn, e := pgPool.Acquire(context.Background())
	if e != nil {
		return nil, e
	}
	defer conn.Release()

	sqlCreateStmt := `
		CREATE TABLE if not exists "users" (
			"id" BIGSERIAL PRIMARY KEY,
			"login" varchar UNIQUE,
			"password" varchar NOT NULL,
			"authtoken" varchar,
			"current" numeric(7,2),
			"withdrawn" numeric(7,2)
		  );

		  CREATE TABLE if not exists "orders" (
			"id" BIGSERIAL PRIMARY KEY,
			"user_id" int NOT NULL,
			"number" varchar NOT NULL,
			"status" varchar,
			"accrual" numeric(7,2),
			"uploaded_at" timestamp DEFAULT (now())
		  );

		  CREATE TABLE if not exists "withdrawals" (
			"id" BIGSERIAL PRIMARY KEY,
			"user_id" int NOT NULL,
			"order" varchar NOT NULL,
			"sum" numeric(7,2),
			"processed_at" timestamp DEFAULT (now())
		  );

		  CREATE INDEX if not exists "order_status" ON "orders" ("user_id", "number", "status");

		  ALTER TABLE if not exists "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

		  ALTER TABLE if not exists "withdrawals" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

		`

	_, err = conn.Exec(context.Background(), sqlCreateStmt)
	if err != nil {
		return nil, err
	}

	return pgPool, nil
}
