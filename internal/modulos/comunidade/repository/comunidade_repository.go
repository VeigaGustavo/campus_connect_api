package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	comunidadeService "campus_connect_api/internal/modulos/comunidade/service"
)

type comunidadeRepositoryPostgres struct {
	store *database.Postgres
}

func NovoComunidadeRepository(store *database.Postgres) comunidadeService.ComunidadeRepository {
	return &comunidadeRepositoryPostgres{store: store}
}

func (repositorio *comunidadeRepositoryPostgres) ListarComunidades(contexto context.Context) ([]comunidadeService.Comunidade, error) {
	return repositorio.store.ListarComunidades(contexto)
}

func (repositorio *comunidadeRepositoryPostgres) InserirComunidade(contexto context.Context, criadoPor string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.store.InserirComunidade(contexto, criadoPor, corpo)
}

func (repositorio *comunidadeRepositoryPostgres) AtualizarComunidade(contexto context.Context, id, usuarioID string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.store.AtualizarComunidade(contexto, id, usuarioID, corpo)
}

func (repositorio *comunidadeRepositoryPostgres) AtualizarComunidadeComoAdmin(contexto context.Context, id string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.store.AtualizarComunidadeComoAdmin(contexto, id, corpo)
}

func (repositorio *comunidadeRepositoryPostgres) RemoverComunidade(contexto context.Context, id, usuarioID string) error {
	return repositorio.store.RemoverComunidade(contexto, id, usuarioID)
}

func (repositorio *comunidadeRepositoryPostgres) RemoverComunidadeComoAdmin(contexto context.Context, id string) error {
	return repositorio.store.RemoverComunidadeComoAdmin(contexto, id)
}
