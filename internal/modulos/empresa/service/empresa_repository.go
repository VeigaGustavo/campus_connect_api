package service

import (
	"context"
)

type EmpresaRepository interface {
	ListarOportunidades(contexto context.Context) ([]Oportunidade, error)
	ObterOportunidade(contexto context.Context, id string) (Oportunidade, bool, error)
	InserirOportunidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarOportunidade) (Oportunidade, error)
	AtualizarOportunidade(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarOportunidade) (Oportunidade, error)
	RemoverOportunidade(contexto context.Context, id, usuarioID, perfil string) error
}
