package service

import (
	"context"
)

type LeituraRepository interface {
	ListarLeituraSemanal(contexto context.Context) ([]ItemLeituraSemanal, error)
	InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error)
	AtualizarLeituraSemanal(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error)
	RemoverLeituraSemanal(contexto context.Context, id, usuarioID, perfil string) error
}
