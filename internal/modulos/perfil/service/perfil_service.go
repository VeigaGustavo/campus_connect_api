package service

import (
	"context"
)

type PerfilService struct {
	repositorio PerfilRepository
}

func NovoPerfilService(repositorio PerfilRepository) *PerfilService {
	return &PerfilService{repositorio: repositorio}
}

func (servico *PerfilService) ObterPerfil(contexto context.Context, usuarioID string) (PerfilUsuario, error) {
	return servico.repositorio.PerfilUsuario(contexto, usuarioID)
}

func (servico *PerfilService) AtualizarPerfil(contexto context.Context, usuarioID string, corpo RequisicaoAtualizarPerfil) (PerfilUsuario, error) {
	return servico.repositorio.AtualizarPerfilUsuario(contexto, usuarioID, corpo)
}

func (servico *PerfilService) HistoricoPerfil(contexto context.Context, usuarioID string, limite int) (RespostaHistoricoPerfil, error) {
	if limite <= 0 || limite > 50 {
		limite = 20
	}
	return servico.repositorio.HistoricoPerfilUsuario(contexto, usuarioID, limite)
}
