package database

import "github.com/jackc/pgx/v5/pgxpool"

// Postgres implementa Armazenamento com PostgreSQL.
type Postgres struct {
	pool *pgxpool.Pool
}

func NovoPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}
