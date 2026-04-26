package service

import (
	"context"
)

type EventoRepository interface {
	ListarEventos(contexto context.Context) ([]EventoCampus, error)
	InserirEvento(contexto context.Context, criadoPor string, corpo RequisicaoEvento) (EventoCampus, error)
	AtualizarEvento(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoEvento) (EventoCampus, error)
	RemoverEvento(contexto context.Context, id, usuarioID, perfil string) error
}
