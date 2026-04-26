package model

type RequisicaoProjeto struct {
	Titulo            string `json:"title"`
	Descricao         string `json:"summary"`
	EscopoPublicacao  string `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string `json:"publish_group_id,omitempty"`
}
