package repository

import (
	"context"

	"campus_connect_api/internal/infra/database"
	empresaService "campus_connect_api/internal/modulos/empresa/service"
)

type empresaRepositoryPostgres struct {
	store *database.Postgres
}

func NovoEmpresaRepository(store *database.Postgres) empresaService.EmpresaRepository {
	return &empresaRepositoryPostgres{store: store}
}

func (repositorio *empresaRepositoryPostgres) ListarOportunidades(contexto context.Context) ([]empresaService.Oportunidade, error) {
	return repositorio.store.ListarOportunidades(contexto)
}

func (repositorio *empresaRepositoryPostgres) ObterOportunidade(contexto context.Context, id string) (empresaService.Oportunidade, bool, error) {
	return repositorio.store.ObterOportunidade(contexto, id)
}

func (repositorio *empresaRepositoryPostgres) InserirOportunidade(contexto context.Context, criadoPor string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	return repositorio.store.InserirOportunidade(contexto, criadoPor, corpo)
}

func (repositorio *empresaRepositoryPostgres) AtualizarOportunidade(contexto context.Context, id, usuarioID, perfil string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	return repositorio.store.AtualizarOportunidade(contexto, id, usuarioID, perfil, corpo)
}

func (repositorio *empresaRepositoryPostgres) RemoverOportunidade(contexto context.Context, id, usuarioID, perfil string) error {
	return repositorio.store.RemoverOportunidade(contexto, id, usuarioID, perfil)
}
