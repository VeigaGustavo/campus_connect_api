package model

type ItemLeituraSemanal struct {
	Identificador string          `json:"id"`
	Tipo          TipoItemLeitura `json:"kind"`
	Titulo        string          `json:"title"`
	Fonte         string          `json:"source"`
	Resumo        string          `json:"excerpt"`
	URLImagem     string          `json:"image_url"`
	RotuloMeta    string          `json:"meta_label"`
}
