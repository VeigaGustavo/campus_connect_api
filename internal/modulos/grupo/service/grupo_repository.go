package service

import (
	"context"
)

type GrupoRepository interface {
	ListarGrupos(contexto context.Context) ([]GrupoEstudo, error)
	InserirGrupo(contexto context.Context, criadoPor string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error)
	AtualizarGrupo(contexto context.Context, id, usuarioID string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error)
	AtualizarGrupoComoAdmin(contexto context.Context, id string, corpo RequisicaoCriarGrupo) (GrupoEstudo, error)
	RemoverGrupo(contexto context.Context, id, usuarioID string) error
	RemoverGrupoComoAdmin(contexto context.Context, id string) error
	ListarMensagensGrupo(contexto context.Context, grupoID string) ([]MensagemChatGrupo, error)
	AdicionarMensagemGrupo(contexto context.Context, grupoID, autorID, texto string) (MensagemChatGrupo, error)
	ListarArquivosGrupo(grupoID string) []ArquivoGrupo
	AdicionarArquivoGrupo(grupoID, autorID, nome, url string) ArquivoGrupo
	ListarReunioesGrupo(grupoID string) []ReuniaoGrupo
	AdicionarReuniaoGrupo(grupoID string, corpo RequisicaoReuniaoGrupo) ReuniaoGrupo
	ListarEventosAssociadosGrupo(grupoID string) []AssociacaoGrupoEvento
	AssociarEventoGrupo(grupoID, eventoID string) AssociacaoGrupoEvento
	ListarLeiturasAssociadasGrupo(grupoID string) []AssociacaoGrupoLeitura
	AssociarLeituraGrupo(grupoID, leituraID string) AssociacaoGrupoLeitura
	ObterVisibilidadeGrupo(contexto context.Context, grupoID string) (string, bool, error)
	InserirMembroGrupo(contexto context.Context, grupoID, usuarioID, papel string) error
	CriarPedidoEntradaGrupo(contexto context.Context, grupoID, usuarioID string) error
	ListarMembrosGrupo(contexto context.Context, grupoID string) ([]MembroGrupo, error)
}
