package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type perfilRepositoryPostgres struct {
	store database.PersistenciaFeedPerfil
}

func NovoPerfilRepository(pool *pgxpool.Pool) perfilService.PerfilRepository {
	return &perfilRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *perfilRepositoryPostgres) PerfilUsuario(contexto context.Context, usuarioID string) (perfilService.PerfilUsuario, error) {
	return repositorio.store.PerfilUsuario(contexto, usuarioID)
}
