package service

import (
	"context"
)

type LeituraRepository interface {
	ListarLeituraSemanal(contexto context.Context, filtroKind string) ([]ItemLeituraSemanal, error)
	InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error)
	AtualizarLeituraSemanal(contexto context.Context, id, usuarioID string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error)
	AtualizarLeituraSemanalComoAdmin(contexto context.Context, id string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error)
	RemoverLeituraSemanal(contexto context.Context, id, usuarioID string) error
	RemoverLeituraSemanalComoAdmin(contexto context.Context, id string) error
}
