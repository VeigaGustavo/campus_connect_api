package service

import (
	"context"
)

type PerfilRepository interface {
	PerfilUsuario(contexto context.Context, usuarioID string) (PerfilUsuario, error)
	AtualizarPerfilUsuario(contexto context.Context, usuarioID string, corpo RequisicaoAtualizarPerfil) (PerfilUsuario, error)
	HistoricoPerfilUsuario(contexto context.Context, usuarioID string, limite int) (RespostaHistoricoPerfil, error)
}
