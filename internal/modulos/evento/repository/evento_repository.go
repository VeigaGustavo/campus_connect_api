package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	eventoService "campus_connect_api/internal/modulos/evento/service"
)

type eventoRepositoryPostgres struct {
	store *database.Postgres
}

func NovoEventoRepository(store *database.Postgres) eventoService.EventoRepository {
	return &eventoRepositoryPostgres{store: store}
}

func (repositorio *eventoRepositoryPostgres) ListarEventos(contexto context.Context) ([]eventoService.EventoCampus, error) {
	return repositorio.store.ListarEventos(contexto)
}

func (repositorio *eventoRepositoryPostgres) InserirEvento(contexto context.Context, criadoPor string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.store.InserirEvento(contexto, criadoPor, corpo)
}

func (repositorio *eventoRepositoryPostgres) AtualizarEvento(contexto context.Context, id, usuarioID string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.store.AtualizarEvento(contexto, id, usuarioID, corpo)
}

func (repositorio *eventoRepositoryPostgres) AtualizarEventoComoAdmin(contexto context.Context, id string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.store.AtualizarEventoComoAdmin(contexto, id, corpo)
}

func (repositorio *eventoRepositoryPostgres) RemoverEvento(contexto context.Context, id, usuarioID string) error {
	return repositorio.store.RemoverEvento(contexto, id, usuarioID)
}

func (repositorio *eventoRepositoryPostgres) RemoverEventoComoAdmin(contexto context.Context, id string) error {
	return repositorio.store.RemoverEventoComoAdmin(contexto, id)
}
