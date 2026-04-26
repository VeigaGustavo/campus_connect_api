package service

import (
	"context"
)

type UniversidadeRepository interface {
	ListarAvisosUniversidade(contexto context.Context) ([]AvisoUniversidade, error)
	InserirAvisoUniversidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error)
	AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarAvisoUniversidade) (AvisoUniversidade, error)
	RemoverAvisoUniversidade(contexto context.Context, id, usuarioID, perfil string) error
}
