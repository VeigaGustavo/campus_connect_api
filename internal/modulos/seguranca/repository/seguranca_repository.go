package repository

import (
	"context"
	"errors"

	"campus_connect_api/internal/infra/database"
	comum "campus_connect_api/internal/modulos/comum"
	segurancaService "campus_connect_api/internal/modulos/seguranca/service"
)

type segurancaRepositoryPostgres struct {
	store *database.Postgres
}

func NovoSegurancaRepository(store *database.Postgres) segurancaService.SegurancaRepository {
	return &segurancaRepositoryPostgres{store: store}
}

func (repositorio *segurancaRepositoryPostgres) Autenticar(contexto context.Context, email, senha string) (*segurancaService.UsuarioAutenticado, error) {
	usuario, err := repositorio.store.Autenticar(contexto, email, senha)
	if err != nil {
		if errors.Is(err, comum.ErrNaoEncontrado) {
			return nil, segurancaService.ErrAutenticacaoInvalida
		}
		return nil, err
	}
	return mapearUsuarioAutenticado(usuario), nil
}

func mapearUsuarioAutenticado(usuario *database.UsuarioInterno) *segurancaService.UsuarioAutenticado {
	if usuario == nil {
		return nil
	}
	return &segurancaService.UsuarioAutenticado{
		ID:           usuario.ID,
		Nome:         usuario.Nome,
		Email:        usuario.Email,
		PerfilCodigo: usuario.PerfilCodigo,
	}
}
