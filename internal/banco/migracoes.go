package banco

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AplicarMigracoesEssenciais aplica ALTERs idempotentes necessários ao código atual.
func AplicarMigracoesEssenciais(contexto context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(contexto, `ALTER TABLE feed_posts ADD COLUMN IF NOT EXISTS content_kind TEXT NOT NULL DEFAULT ''`); err != nil {
		return err
	}
	_, err := pool.Exec(contexto, `
ALTER TABLE feed_cartoes DROP CONSTRAINT IF EXISTS feed_cartoes_kind_check;
ALTER TABLE feed_cartoes ADD CONSTRAINT feed_cartoes_kind_check CHECK (kind::text = ANY (ARRAY[
  'internship','event','study_group','project','notice','campus_feed','post','reading'
]::text[]));`)
	return err
}
