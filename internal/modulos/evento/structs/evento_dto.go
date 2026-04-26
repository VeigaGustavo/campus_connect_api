package model

type RequisicaoEvento struct {
	Titulo            string `json:"title"`
	Descricao         string `json:"description"`
	InicioEm          string `json:"start_at"`
	Local             string `json:"location"`
	Organizador       string `json:"organizer"`
	EscopoPublicacao  string `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string `json:"publish_group_id,omitempty"`
}
