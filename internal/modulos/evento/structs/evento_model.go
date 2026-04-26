package model

import comum "campus_connect_api/internal/modulos/comum"

type EventoCampus struct {
	Identificador string                   `json:"id"`
	AutorID       string                   `json:"author_id"`
	Autor         comum.PerfilPublicoAutor `json:"author"`
	Titulo        string                   `json:"title"`
	Descricao     string                   `json:"description"`
	InicioEm      string                   `json:"start_at"`
	Local         string                   `json:"location"`
	Organizador   string                   `json:"organizer"`
}
