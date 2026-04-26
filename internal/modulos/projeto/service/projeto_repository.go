package service

import (
	"context"
)

type ProjetoRepository interface {
	ListarProjetos(contexto context.Context) ([]Projeto, error)
	InserirProjeto(contexto context.Context, criadoPor string, corpo RequisicaoProjeto) (Projeto, error)
	AtualizarProjeto(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoProjeto) (Projeto, error)
	RemoverProjeto(contexto context.Context, id, usuarioID, perfil string) error
}
