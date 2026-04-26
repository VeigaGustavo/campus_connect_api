package service

import (
	"context"
)

type EventoService struct {
	repositorio EventoRepository
}

func NovoEventoService(repositorio EventoRepository) *EventoService {
	return &EventoService{repositorio: repositorio}
}

func (servico *EventoService) ListarEventos(contexto context.Context) ([]EventoCampus, error) {
	return servico.repositorio.ListarEventos(contexto)
}

func (servico *EventoService) CriarEvento(contexto context.Context, criadoPor string, corpo RequisicaoEvento) (EventoCampus, error) {
	return servico.repositorio.InserirEvento(contexto, criadoPor, corpo)
}

func (servico *EventoService) AtualizarEvento(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoEvento) (EventoCampus, error) {
	if perfil == "sistema_admin" {
		return servico.repositorio.AtualizarEventoComoAdmin(contexto, id, corpo)
	}
	return servico.repositorio.AtualizarEvento(contexto, id, usuarioID, corpo)
}

func (servico *EventoService) RemoverEvento(contexto context.Context, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return servico.repositorio.RemoverEventoComoAdmin(contexto, id)
	}
	return servico.repositorio.RemoverEvento(contexto, id, usuarioID)
}
