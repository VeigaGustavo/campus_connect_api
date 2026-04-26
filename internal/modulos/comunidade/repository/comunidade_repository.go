package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	comunidadeService "campus_connect_api/internal/modulos/comunidade/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type comunidadeRepositoryPostgres struct {
	store database.PersistenciaComunidade
}

func NovoComunidadeRepository(pool *pgxpool.Pool) comunidadeService.ComunidadeRepository {
	return &comunidadeRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *comunidadeRepositoryPostgres) ListarComunidades(contexto context.Context) ([]comunidadeService.Comunidade, error) {
	return repositorio.store.ListarComunidades(contexto)
}

func (repositorio *comunidadeRepositoryPostgres) InserirComunidade(contexto context.Context, criadoPor string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.store.InserirComunidade(contexto, criadoPor, corpo)
}

func (repositorio *comunidadeRepositoryPostgres) AtualizarComunidade(contexto context.Context, id, usuarioID, perfil string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.store.AtualizarComunidade(contexto, id, usuarioID, perfil, corpo)
}

func (repositorio *comunidadeRepositoryPostgres) RemoverComunidade(contexto context.Context, id, usuarioID, perfil string) error {
	return repositorio.store.RemoverComunidade(contexto, id, usuarioID, perfil)
}
