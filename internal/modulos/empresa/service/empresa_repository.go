package service

import (
	"context"
)

type EmpresaRepository interface {
	ListarOportunidades(contexto context.Context) ([]Oportunidade, error)
	ObterOportunidade(contexto context.Context, id string) (Oportunidade, bool, error)
	InserirOportunidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarOportunidade) (Oportunidade, error)
	AtualizarOportunidade(contexto context.Context, id, usuarioID string, corpo RequisicaoCriarOportunidade) (Oportunidade, error)
	AtualizarOportunidadeComoAdmin(contexto context.Context, id string, corpo RequisicaoCriarOportunidade) (Oportunidade, error)
	RemoverOportunidade(contexto context.Context, id, usuarioID string) error
	RemoverOportunidadeComoAdmin(contexto context.Context, id string) error
}
