package service

import (
	"context"
)

type GrupoService struct {
	repositorio GrupoRepository
}

func NovoGrupoService(repositorio GrupoRepository) *GrupoService {
	return &GrupoService{repositorio: repositorio}
}

func (servico *GrupoService) ListarGrupos(contexto context.Context) ([]GrupoEstudo, error) {
	return servico.repositorio.ListarGrupos(contexto)
}

func (servico *GrupoService) CriarGrupo(contexto context.Context, criadoPor string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error) {
	return servico.repositorio.InserirGrupo(contexto, criadoPor, corpo)
}

func (servico *GrupoService) AtualizarGrupo(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error) {
	return servico.repositorio.AtualizarGrupo(contexto, id, usuarioID, perfil, corpo)
}

func (servico *GrupoService) RemoverGrupo(contexto context.Context, id, usuarioID, perfil string) error {
	return servico.repositorio.RemoverGrupo(contexto, id, usuarioID, perfil)
}

func (servico *GrupoService) ListarMensagensGrupo(grupoID string) []MensagemChatGrupo {
	return servico.repositorio.ListarMensagensGrupo(grupoID)
}

func (servico *GrupoService) AdicionarMensagemGrupo(grupoID, autorID, texto string) MensagemChatGrupo {
	return servico.repositorio.AdicionarMensagemGrupo(grupoID, autorID, texto)
}

func (servico *GrupoService) ListarArquivosGrupo(grupoID string) []ArquivoGrupo {
	return servico.repositorio.ListarArquivosGrupo(grupoID)
}

func (servico *GrupoService) AdicionarArquivoGrupo(grupoID, autorID, nome, url string) ArquivoGrupo {
	return servico.repositorio.AdicionarArquivoGrupo(grupoID, autorID, nome, url)
}

func (servico *GrupoService) ListarReunioesGrupo(grupoID string) []ReuniaoGrupo {
	return servico.repositorio.ListarReunioesGrupo(grupoID)
}

func (servico *GrupoService) AdicionarReuniaoGrupo(grupoID string, corpo RequisicaoReuniaoGrupo) ReuniaoGrupo {
	return servico.repositorio.AdicionarReuniaoGrupo(grupoID, corpo)
}

func (servico *GrupoService) ListarEventosAssociadosGrupo(grupoID string) []AssociacaoGrupoEvento {
	return servico.repositorio.ListarEventosAssociadosGrupo(grupoID)
}

func (servico *GrupoService) AssociarEventoGrupo(grupoID, eventoID string) AssociacaoGrupoEvento {
	return servico.repositorio.AssociarEventoGrupo(grupoID, eventoID)
}

func (servico *GrupoService) ListarLeiturasAssociadasGrupo(grupoID string) []AssociacaoGrupoLeitura {
	return servico.repositorio.ListarLeiturasAssociadasGrupo(grupoID)
}

func (servico *GrupoService) AssociarLeituraGrupo(grupoID, leituraID string) AssociacaoGrupoLeitura {
	return servico.repositorio.AssociarLeituraGrupo(grupoID, leituraID)
}
