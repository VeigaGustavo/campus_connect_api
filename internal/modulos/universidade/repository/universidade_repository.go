package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	universidadeService "campus_connect_api/internal/modulos/universidade/service"
)

type universidadeRepositoryPostgres struct {
	store *database.Postgres
}

func NovoUniversidadeRepository(store *database.Postgres) universidadeService.UniversidadeRepository {
	return &universidadeRepositoryPostgres{store: store}
}

func (repositorio *universidadeRepositoryPostgres) ListarAvisosUniversidade(contexto context.Context) ([]universidadeService.AvisoUniversidade, error) {
	return repositorio.store.ListarAvisosUniversidade(contexto)
}

func (repositorio *universidadeRepositoryPostgres) InserirAvisoUniversidade(contexto context.Context, criadoPor string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return repositorio.store.InserirAvisoUniversidade(contexto, criadoPor, corpo)
}

func (repositorio *universidadeRepositoryPostgres) AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return repositorio.store.AtualizarAvisoUniversidade(contexto, id, usuarioID, corpo)
}

func (repositorio *universidadeRepositoryPostgres) AtualizarAvisoUniversidadeComoAdmin(contexto context.Context, id string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return repositorio.store.AtualizarAvisoUniversidadeComoAdmin(contexto, id, corpo)
}

func (repositorio *universidadeRepositoryPostgres) RemoverAvisoUniversidade(contexto context.Context, id, usuarioID string) error {
	return repositorio.store.RemoverAvisoUniversidade(contexto, id, usuarioID)
}

func (repositorio *universidadeRepositoryPostgres) RemoverAvisoUniversidadeComoAdmin(contexto context.Context, id string) error {
	return repositorio.store.RemoverAvisoUniversidadeComoAdmin(contexto, id)
}
