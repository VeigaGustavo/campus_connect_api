package service

import (
	"context"
)

type ProjetoService struct {
	repositorio ProjetoRepository
}

func NovoProjetoService(repositorio ProjetoRepository) *ProjetoService {
	return &ProjetoService{repositorio: repositorio}
}

func (servico *ProjetoService) ListarProjetos(contexto context.Context) ([]Projeto, error) {
	return servico.repositorio.ListarProjetos(contexto)
}

func (servico *ProjetoService) CriarProjeto(contexto context.Context, criadoPor string, corpo RequisicaoProjeto) (Projeto, error) {
	return servico.repositorio.InserirProjeto(contexto, criadoPor, corpo)
}

func (servico *ProjetoService) AtualizarProjeto(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoProjeto) (Projeto, error) {
	if perfil == "sistema_admin" {
		return servico.repositorio.AtualizarProjetoComoAdmin(contexto, id, corpo)
	}
	return servico.repositorio.AtualizarProjeto(contexto, id, usuarioID, corpo)
}

func (servico *ProjetoService) RemoverProjeto(contexto context.Context, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return servico.repositorio.RemoverProjetoComoAdmin(contexto, id)
	}
	return servico.repositorio.RemoverProjeto(contexto, id, usuarioID)
}
