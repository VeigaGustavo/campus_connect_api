package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	projetoService "campus_connect_api/internal/modulos/projeto/service"
)

type projetoRepositoryPostgres struct {
	store *database.Postgres
}

func NovoProjetoRepository(store *database.Postgres) projetoService.ProjetoRepository {
	return &projetoRepositoryPostgres{store: store}
}

func (repositorio *projetoRepositoryPostgres) ListarProjetos(contexto context.Context) ([]projetoService.Projeto, error) {
	return repositorio.store.ListarProjetos(contexto)
}

func (repositorio *projetoRepositoryPostgres) InserirProjeto(contexto context.Context, criadoPor string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return repositorio.store.InserirProjeto(contexto, criadoPor, corpo)
}

func (repositorio *projetoRepositoryPostgres) AtualizarProjeto(contexto context.Context, id, usuarioID, perfil string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return repositorio.store.AtualizarProjeto(contexto, id, usuarioID, perfil, corpo)
}

func (repositorio *projetoRepositoryPostgres) RemoverProjeto(contexto context.Context, id, usuarioID, perfil string) error {
	return repositorio.store.RemoverProjeto(contexto, id, usuarioID, perfil)
}
