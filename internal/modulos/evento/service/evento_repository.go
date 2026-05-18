package service

import (
	"context"
)

type EventoRepository interface {
	ListarEventos(contexto context.Context) ([]EventoCampus, error)
	InserirEvento(contexto context.Context, criadoPor string, corpo RequisicaoEvento) (EventoCampus, error)
	AtualizarEvento(contexto context.Context, id, usuarioID string, corpo RequisicaoEvento) (EventoCampus, error)
	AtualizarEventoComoAdmin(contexto context.Context, id string, corpo RequisicaoEvento) (EventoCampus, error)
	RemoverEvento(contexto context.Context, id, usuarioID string) error
	RemoverEventoComoAdmin(contexto context.Context, id string) error
}
