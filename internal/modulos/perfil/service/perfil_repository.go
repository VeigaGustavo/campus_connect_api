package service

import (
	"context"
)

type PerfilRepository interface {
	PerfilUsuario(contexto context.Context, usuarioID string) (PerfilUsuario, error)
	PerfilParaExibicao(contexto context.Context, usuarioID, perfilCodigoConta string) (PerfilUsuario, error)
	AtualizarPerfilUsuario(contexto context.Context, usuarioID string, corpo RequisicaoAtualizarPerfil) (PerfilUsuario, error)
	AtualizarPerfilParaExibicao(contexto context.Context, usuarioID, perfilCodigoConta string, corpo RequisicaoAtualizarPerfil) (PerfilUsuario, error)
	HistoricoPerfilUsuario(contexto context.Context, usuarioID string, limite int) (RespostaHistoricoPerfil, error)
	AtualizarURLImagemPerfil(contexto context.Context, usuarioID, perfilCodigoConta, tipoImagem, urlPublica string) error
}
