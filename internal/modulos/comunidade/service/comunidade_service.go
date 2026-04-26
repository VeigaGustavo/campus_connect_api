package service

import (
	"context"
)

type ComunidadeService struct {
	repositorio ComunidadeRepository
}

func NovoComunidadeService(repositorio ComunidadeRepository) *ComunidadeService {
	return &ComunidadeService{repositorio: repositorio}
}

func (servico *ComunidadeService) ListarComunidades(contexto context.Context) ([]Comunidade, error) {
	return servico.repositorio.ListarComunidades(contexto)
}

func (servico *ComunidadeService) CriarComunidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarComunidade) (Comunidade, error) {
	return servico.repositorio.InserirComunidade(contexto, criadoPor, corpo)
}

func (servico *ComunidadeService) AtualizarComunidade(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarComunidade) (Comunidade, error) {
	return servico.repositorio.AtualizarComunidade(contexto, id, usuarioID, perfil, corpo)
}

func (servico *ComunidadeService) RemoverComunidade(contexto context.Context, id, usuarioID, perfil string) error {
	return servico.repositorio.RemoverComunidade(contexto, id, usuarioID, perfil)
}
