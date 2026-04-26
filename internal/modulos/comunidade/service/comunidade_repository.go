package service

import (
	"context"
)

type ComunidadeRepository interface {
	ListarComunidades(contexto context.Context) ([]Comunidade, error)
	InserirComunidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarComunidade) (Comunidade, error)
	AtualizarComunidade(contexto context.Context, id, usuarioID string, corpo RequisicaoCriarComunidade) (Comunidade, error)
	AtualizarComunidadeComoAdmin(contexto context.Context, id string, corpo RequisicaoCriarComunidade) (Comunidade, error)
	RemoverComunidade(contexto context.Context, id, usuarioID string) error
	RemoverComunidadeComoAdmin(contexto context.Context, id string) error
}
