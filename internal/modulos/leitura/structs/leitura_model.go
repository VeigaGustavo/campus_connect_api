package model

import comum "campus_connect_api/internal/modulos/comum"

type ItemLeituraSemanal struct {
	Identificador string                   `json:"id"`
	AutorID       string                   `json:"author_id"`
	Autor         comum.PerfilPublicoAutor `json:"author"`
	Tipo          TipoItemLeitura          `json:"kind"`
	Titulo        string                   `json:"title"`
	Fonte         string                   `json:"source"`
	Resumo        string                   `json:"excerpt"`
	URLImagem     string                   `json:"image_url"`
	RotuloMeta    string                   `json:"meta_label"`
}

type RespostaListaLeituraSemanal struct {
	Itens []ItemLeituraSemanal `json:"items"`
}
