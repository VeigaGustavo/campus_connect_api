package service

import (
	"context"
)

type ProjetoRepository interface {
	ListarProjetos(contexto context.Context) ([]Projeto, error)
	InserirProjeto(contexto context.Context, criadoPor string, corpo RequisicaoProjeto) (Projeto, error)
	AtualizarProjeto(contexto context.Context, id, usuarioID string, corpo RequisicaoProjeto) (Projeto, error)
	AtualizarProjetoComoAdmin(contexto context.Context, id string, corpo RequisicaoProjeto) (Projeto, error)
	RemoverProjeto(contexto context.Context, id, usuarioID string) error
	RemoverProjetoComoAdmin(contexto context.Context, id string) error
}
