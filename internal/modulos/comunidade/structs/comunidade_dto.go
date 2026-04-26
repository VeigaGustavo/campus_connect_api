package model

type RequisicaoCriarComunidade struct {
	Nome      string `json:"name"`
	Tipo      string `json:"kind"`
	Descricao string `json:"description"`
}
