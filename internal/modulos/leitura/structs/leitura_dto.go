package model

type RequisicaoLeituraSemanal struct {
	Tipo              string `json:"kind"`
	Titulo            string `json:"title"`
	Fonte             string `json:"source"`
	Resumo            string `json:"excerpt"`
	URLImagem         string `json:"image_url"`
	RotuloMeta        string `json:"meta_label"`
	EscopoPublicacao  string `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string `json:"publish_group_id,omitempty"`
}
