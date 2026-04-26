package service

import (
	"context"
)

type LeituraService struct {
	repositorio LeituraRepository
}

func NovoLeituraService(repositorio LeituraRepository) *LeituraService {
	return &LeituraService{repositorio: repositorio}
}

func (servico *LeituraService) ListarLeituraSemanal(contexto context.Context) ([]ItemLeituraSemanal, error) {
	return servico.repositorio.ListarLeituraSemanal(contexto)
}

func (servico *LeituraService) CriarLeituraSemanal(contexto context.Context, criadoPor string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error) {
	return servico.repositorio.InserirLeituraSemanal(contexto, criadoPor, corpo)
}

func (servico *LeituraService) AtualizarLeituraSemanal(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoLeituraSemanal) (ItemLeituraSemanal, error) {
	return servico.repositorio.AtualizarLeituraSemanal(contexto, id, usuarioID, perfil, corpo)
}

func (servico *LeituraService) RemoverLeituraSemanal(contexto context.Context, id, usuarioID, perfil string) error {
	return servico.repositorio.RemoverLeituraSemanal(contexto, id, usuarioID, perfil)
}
