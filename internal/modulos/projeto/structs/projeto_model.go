package model

type Projeto struct {
	Identificador string `json:"id"`
	Titulo        string `json:"title"`
	Descricao     string `json:"summary"`
}
