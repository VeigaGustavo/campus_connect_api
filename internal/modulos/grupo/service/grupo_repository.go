package service

import (
	"context"
)

type GrupoRepository interface {
	ListarGrupos(contexto context.Context) ([]GrupoEstudo, error)
	InserirGrupo(contexto context.Context, criadoPor string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error)
	AtualizarGrupo(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error)
	RemoverGrupo(contexto context.Context, id, usuarioID, perfil string) error
	ListarMensagensGrupo(grupoID string) []MensagemChatGrupo
	AdicionarMensagemGrupo(grupoID, autorID, texto string) MensagemChatGrupo
	ListarArquivosGrupo(grupoID string) []ArquivoGrupo
	AdicionarArquivoGrupo(grupoID, autorID, nome, url string) ArquivoGrupo
	ListarReunioesGrupo(grupoID string) []ReuniaoGrupo
	AdicionarReuniaoGrupo(grupoID string, corpo RequisicaoReuniaoGrupo) ReuniaoGrupo
	ListarEventosAssociadosGrupo(grupoID string) []AssociacaoGrupoEvento
	AssociarEventoGrupo(grupoID, eventoID string) AssociacaoGrupoEvento
	ListarLeiturasAssociadasGrupo(grupoID string) []AssociacaoGrupoLeitura
	AssociarLeituraGrupo(grupoID, leituraID string) AssociacaoGrupoLeitura
}
