package service

import (
	"context"

	comum "campus_connect_api/internal/modulos/comum"
)

type DiscoverService struct {
	repositorio DiscoverRepository
}

func NovoDiscoverService(repositorio DiscoverRepository) *DiscoverService {
	return &DiscoverService{repositorio: repositorio}
}

func (servico *DiscoverService) FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) (RespostaDescobrir, error) {
	if filtro == "" {
		filtro = comum.FiltroDescobrirTodos
	}
	itens, err := servico.repositorio.FeedDescobrir(contexto, filtro, gruposDoUsuario)
	if err != nil {
		return RespostaDescobrir{}, err
	}
	return RespostaDescobrir{Itens: itens}, nil
}
