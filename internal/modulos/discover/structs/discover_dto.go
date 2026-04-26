package model

type ItemDescobrir struct {
	Identificador     string                 `json:"id"`
	Categoria         CategoriaItemDescobrir `json:"kind"`
	Titulo            string                 `json:"title"`
	Subtitulo         string                 `json:"subtitle"`
	Resumo            string                 `json:"excerpt"`
	MetaPrincipal     string                 `json:"meta_primary"`
	MetaSecundaria    string                 `json:"meta_secondary"`
	IDReferencia      string                 `json:"reference_id"`
	EscopoPublicacao  string                 `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string                 `json:"publish_group_id,omitempty"`
}

type RespostaDescobrir struct {
	Itens []ItemDescobrir `json:"items"`
}
