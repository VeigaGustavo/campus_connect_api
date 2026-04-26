package model

type Comunidade struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	Tipo          string `json:"kind"`
	Descricao     string `json:"description"`
}
