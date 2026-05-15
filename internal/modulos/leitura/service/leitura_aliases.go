package service

import model "campus_connect_api/internal/modulos/leitura/structs"

type TipoItemLeitura = model.TipoItemLeitura

const (
	LeituraNoticiaCampus = model.LeituraNoticiaCampus
	LeituraRevista       = model.LeituraRevista
	LeituraArtigo        = model.LeituraArtigo
)

type ItemLeituraSemanal = model.ItemLeituraSemanal
type RespostaListaLeituraSemanal = model.RespostaListaLeituraSemanal
type RequisicaoLeituraSemanal = model.RequisicaoLeituraSemanal
