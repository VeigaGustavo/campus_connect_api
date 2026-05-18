package banco

import (
	"context"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"campus_connect_api/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AplicarMigracoes(contexto context.Context, pool *pgxpool.Pool) error {
	entradas, err := fs.ReadDir(db.Migracoes, "migrations")
	if err != nil {
		return fmt.Errorf("ler migracoes: %w", err)
	}
	nomes := make([]string, 0, len(entradas))
	for _, entrada := range entradas {
		if entrada.IsDir() || !strings.HasSuffix(entrada.Name(), ".sql") {
			continue
		}
		nomes = append(nomes, entrada.Name())
	}
	sort.Strings(nomes)

	for _, nome := range nomes {
		conteudo, err := db.Migracoes.ReadFile("migrations/" + nome)
		if err != nil {
			return fmt.Errorf("ler %s: %w", nome, err)
		}
		if _, err := pool.Exec(contexto, string(conteudo)); err != nil {
			return fmt.Errorf("%s: %w", nome, err)
		}
	}
	return nil
}
