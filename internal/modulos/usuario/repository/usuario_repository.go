package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type usuarioRepositoryPostgres struct {
	store database.PersistenciaUsuario
}

func NovoUsuarioRepository(pool *pgxpool.Pool) usuarioService.UsuarioRepository {
	return &usuarioRepositoryPostgres{store: database.NovoPostgres(pool)}
}

func (repositorio *usuarioRepositoryPostgres) CriarUsuario(contexto context.Context, nome, email, senha, perfilCodigo string) (*usuarioService.UsuarioInterno, error) {
	usuario, err := repositorio.store.CriarUsuario(contexto, nome, email, senha, perfilCodigo)
	if err != nil {
		return nil, err
	}
	return mapearUsuarioInterno(usuario), nil
}

func (repositorio *usuarioRepositoryPostgres) CriarUsuarioComCadastro(contexto context.Context, requisicao usuarioService.RequisicaoCadastroUsuario) (*usuarioService.UsuarioInterno, error) {
	usuario, err := repositorio.store.CriarUsuarioComCadastro(contexto, requisicao)
	if err != nil {
		return nil, err
	}
	return mapearUsuarioInterno(usuario), nil
}

func mapearUsuarioInterno(usuario *database.UsuarioInterno) *usuarioService.UsuarioInterno {
	if usuario == nil {
		return nil
	}
	return &usuarioService.UsuarioInterno{
		ID:           usuario.ID,
		Nome:         usuario.Nome,
		Email:        usuario.Email,
		PerfilCodigo: usuario.PerfilCodigo,
	}
}
