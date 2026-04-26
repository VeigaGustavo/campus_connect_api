package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	leituraService "campus_connect_api/internal/modulos/leitura/service"
)

type leituraRepositoryPostgres struct {
	store *database.Postgres
}

func NovoLeituraRepository(store *database.Postgres) leituraService.LeituraRepository {
	return &leituraRepositoryPostgres{store: store}
}

func (repositorio *leituraRepositoryPostgres) ListarLeituraSemanal(contexto context.Context) ([]leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.ListarLeituraSemanal(contexto)
}

func (repositorio *leituraRepositoryPostgres) InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.InserirLeituraSemanal(contexto, criadoPor, corpo)
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanal(contexto context.Context, id, usuarioID string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.AtualizarLeituraSemanal(contexto, id, usuarioID, corpo)
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanalComoAdmin(contexto context.Context, id string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.store.AtualizarLeituraSemanalComoAdmin(contexto, id, corpo)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanal(contexto context.Context, id, usuarioID string) error {
	return repositorio.store.RemoverLeituraSemanal(contexto, id, usuarioID)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanalComoAdmin(contexto context.Context, id string) error {
	return repositorio.store.RemoverLeituraSemanalComoAdmin(contexto, id)
}
