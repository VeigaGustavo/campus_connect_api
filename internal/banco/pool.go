package banco

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NovoPool abre pool de conexões PostgreSQL a partir de DATABASE_URL.
func NovoPool(contexto context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConnLifetime = time.Hour
	cfg.ConnConfig.ConnectTimeout = 10 * time.Second
	return pgxpool.NewWithConfig(contexto, cfg)
}
