package service

import (
	"context"
)

type UniversidadeRepository interface {
	ListarAvisosUniversidade(contexto context.Context) ([]AvisoUniversidade, error)
	InserirAvisoUniversidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error)
	AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error)
	AtualizarAvisoUniversidadeComoAdmin(contexto context.Context, id string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error)
	RemoverAvisoUniversidade(contexto context.Context, id, usuarioID string) error
	RemoverAvisoUniversidadeComoAdmin(contexto context.Context, id string) error
}
