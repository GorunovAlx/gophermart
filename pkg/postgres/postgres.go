package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GorunovAlx/gophermart/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

// Postgres
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

// PGXStdLogger prints pgx logs to the standard logger.
// os.Stderr by default.
type PGXStdLogger struct{}

// New
func New(cfg *config.Config) (*Postgres, error) {
	pgxLogLevel, err := LogLevelFromCfg(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("pgx configuration: %w", err)
	}

	poolConfig.ConnConfig.Logger = &PGXStdLogger{}
	if pgxLogLevel != 0 {
		poolConfig.ConnConfig.LogLevel = pgxLogLevel
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (l *PGXStdLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

// LogLevelFromEnv returns the pgx.LogLevel from the environment variable PGX_LOG_LEVEL.
// By default this is info (pgx.LogLevelInfo), which is good for development.
func LogLevelFromCfg(cfg *config.Config) (pgx.LogLevel, error) {
	if level := cfg.PgxLogLevel; level != "" {
		l, err := pgx.LogLevelFromString(level)
		if err != nil {
			return pgx.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return pgx.LogLevelInfo, nil
}

// Close
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
