package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	eventoService "campus_connect_api/internal/modulos/evento/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type eventoRepositoryPostgres struct {
	store *database.Postgres
}

func NovoEventoRepository(pool *pgxpool.Pool) eventoService.EventoRepository {
	return &eventoRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *eventoRepositoryPostgres) ListarEventos(contexto context.Context) ([]eventoService.EventoCampus, error) {
	return repositorio.store.ListarEventos(contexto)
}

func (repositorio *eventoRepositoryPostgres) InserirEvento(contexto context.Context, criadoPor string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.store.InserirEvento(contexto, criadoPor, corpo)
}

func (repositorio *eventoRepositoryPostgres) AtualizarEvento(contexto context.Context, id, usuarioID, perfil string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.store.AtualizarEvento(contexto, id, usuarioID, perfil, corpo)
}

func (repositorio *eventoRepositoryPostgres) RemoverEvento(contexto context.Context, id, usuarioID, perfil string) error {
	return repositorio.store.RemoverEvento(contexto, id, usuarioID, perfil)
}
