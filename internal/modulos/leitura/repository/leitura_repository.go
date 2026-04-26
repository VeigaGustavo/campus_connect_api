package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	leituraService "campus_connect_api/internal/modulos/leitura/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type leituraRepositoryPostgres struct {
	store database.PersistenciaLeitura
}

func NovoLeituraRepository(pool *pgxpool.Pool) leituraService.LeituraRepository {
	return &leituraRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *leituraRepositoryPostgres) ListarLeituraSemanal(contexto context.Context) ([]leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.ListarLeituraSemanal(contexto)
}

func (repositorio *leituraRepositoryPostgres) InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.InserirLeituraSemanal(contexto, criadoPor, corpo)
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanal(contexto context.Context, id, usuarioID, perfil string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.AtualizarLeituraSemanal(contexto, id, usuarioID, perfil, corpo)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanal(contexto context.Context, id, usuarioID, perfil string) error {
	return repositorio.store.RemoverLeituraSemanal(contexto, id, usuarioID, perfil)
}
