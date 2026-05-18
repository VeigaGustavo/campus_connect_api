package service

import model "campus_connect_api/internal/modulos/grupo/structs"

type NivelGrupoEstudo = model.NivelGrupoEstudo

const (
	NivelIniciante     = model.NivelIniciante
	NivelIntermediario = model.NivelIntermediario
	NivelAvancado      = model.NivelAvancado
)

type GrupoEstudo = model.GrupoEstudo
type AutorGrupo = model.AutorGrupo

var AutorGrupoDe = model.AutorGrupoDe
type MembroGrupo = model.MembroGrupo
type RequisicaoCriarGrupo = model.RequisicaoCriarGrupo
type MensagemChatGrupo = model.MensagemChatGrupo
type RequisicaoMensagemChatGrupo = model.RequisicaoMensagemChatGrupo
type ArquivoGrupo = model.ArquivoGrupo
type RequisicaoArquivoGrupo = model.RequisicaoArquivoGrupo
type ReuniaoGrupo = model.ReuniaoGrupo
type RequisicaoReuniaoGrupo = model.RequisicaoReuniaoGrupo
type AssociacaoGrupoEvento = model.AssociacaoGrupoEvento
type RequisicaoAssociarEventoGrupo = model.RequisicaoAssociarEventoGrupo
type AssociacaoGrupoLeitura = model.AssociacaoGrupoLeitura
type RequisicaoAssociarLeituraGrupo = model.RequisicaoAssociarLeituraGrupo
