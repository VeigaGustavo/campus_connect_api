package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
)

type perfilRepositoryPostgres struct {
	store *database.Postgres
}

func NovoPerfilRepository(store *database.Postgres) perfilService.PerfilRepository {
	return &perfilRepositoryPostgres{store: store}
}

func (repositorio *perfilRepositoryPostgres) PerfilUsuario(contexto context.Context, usuarioID string) (perfilService.PerfilUsuario, error) {
	return repositorio.store.PerfilUsuario(contexto, usuarioID)
}
