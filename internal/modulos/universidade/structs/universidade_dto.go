package model

type RequisicaoCriarAvisoUniversidade struct {
	Titulo            string `json:"title"`
	Descricao         string `json:"description"`
	EscopoPublicacao  string `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string `json:"publish_group_id,omitempty"`
}
