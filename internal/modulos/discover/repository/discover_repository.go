package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	discoverService "campus_connect_api/internal/modulos/discover/service"
)

type discoverRepositoryPostgres struct {
	store *database.Postgres
}

func NovoDiscoverRepository(store *database.Postgres) discoverService.DiscoverRepository {
	return &discoverRepositoryPostgres{store: store}
}

func (repositorio *discoverRepositoryPostgres) FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]discoverService.ItemDescobrir, error) {
	return repositorio.store.FeedDescobrir(contexto, filtro, gruposDoUsuario)
}
