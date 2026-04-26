package service

import (
	"context"
)

type UniversidadeService struct {
	repositorio UniversidadeRepository
}

func NovoUniversidadeService(repositorio UniversidadeRepository) *UniversidadeService {
	return &UniversidadeService{repositorio: repositorio}
}

func (servico *UniversidadeService) ListarAvisosUniversidade(contexto context.Context) ([]AvisoUniversidade, error) {
	return servico.repositorio.ListarAvisosUniversidade(contexto)
}

func (servico *UniversidadeService) CriarAvisoUniversidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error) {
	return servico.repositorio.InserirAvisoUniversidade(contexto, criadoPor, corpo)
}

func (servico *UniversidadeService) AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error) {
	if perfil == "sistema_admin" {
		return servico.repositorio.AtualizarAvisoUniversidadeComoAdmin(contexto, id, corpo)
	}
	return servico.repositorio.AtualizarAvisoUniversidade(contexto, id, usuarioID, corpo)
}

func (servico *UniversidadeService) RemoverAvisoUniversidade(contexto context.Context, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return servico.repositorio.RemoverAvisoUniversidadeComoAdmin(contexto, id)
	}
	return servico.repositorio.RemoverAvisoUniversidade(contexto, id, usuarioID)
}
