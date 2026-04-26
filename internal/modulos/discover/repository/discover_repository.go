package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	discoverService "campus_connect_api/internal/modulos/discover/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type discoverRepositoryPostgres struct {
	store database.PersistenciaFeedPerfil
}

func NovoDiscoverRepository(pool *pgxpool.Pool) discoverService.DiscoverRepository {
	return &discoverRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *discoverRepositoryPostgres) FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]discoverService.ItemDescobrir, error) {
	return repositorio.store.FeedDescobrir(contexto, filtro, gruposDoUsuario)
}
