package model

import comum "campus_connect_api/internal/modulos/comum"

type AvisoUniversidade struct {
	Identificador string                   `json:"id"`
	AutorID       string                   `json:"author_id"`
	Autor         comum.PerfilPublicoAutor `json:"author"`
	Titulo        string                   `json:"title"`
	Descricao     string                   `json:"description"`
}
