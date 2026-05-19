package banco

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AplicarMigracoesEssenciais(contexto context.Context, pool *pgxpool.Pool) error {
	return AplicarMigracoes(contexto, pool)
}
