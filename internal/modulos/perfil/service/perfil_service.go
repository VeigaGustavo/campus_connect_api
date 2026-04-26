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
