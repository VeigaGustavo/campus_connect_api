package model

import comum "campus_connect_api/internal/modulos/comum"

type ItemFeed struct {
	Identificador  string            `json:"id"`
	Categoria      CategoriaItemFeed `json:"kind"`
	Titulo         string            `json:"title"`
	Subtitulo      string            `json:"subtitle"`
	Resumo         string            `json:"excerpt"`
	MetaPrincipal  string            `json:"meta_primary"`
	MetaSecundaria string            `json:"meta_secondary"`
	IDReferencia      string            `json:"reference_id"`
	EscopoPublicacao  string            `json:"publish_scope"`
	IDGrupoPublicacao string            `json:"publish_group_id,omitempty"`
}

type RespostaFeed struct {
	Itens []ItemFeed `json:"items"`
}

type AnexoPost struct {
	Tipo string `json:"type"`
	URL  string `json:"url"`
	Nome string `json:"name,omitempty"`
}

type RequisicaoCriarPost struct {
	Texto             string      `json:"text"`
	Anexos            []AnexoPost `json:"attachments"`
	EscopoPublicacao  string      `json:"publish_scope"`
	IDGrupoPublicacao string      `json:"publish_group_id"`
	TipoConteudo      string      `json:"content_kind,omitempty"`
}

type ComentarioPost struct {
	Identificador  string                   `json:"id"`
	PostID         string                   `json:"post_id"`
	AutorID        string                   `json:"author_id"`
	Autor          comum.PerfilPublicoAutor `json:"author"`
	Texto          string                   `json:"text"`
	GosteiTotal    int                      `json:"likes_count"`
	DesgosteiTotal int                      `json:"dislikes_count"`
	CriadoEm       string                   `json:"created_at"`
}

type RequisicaoCriarComentario struct {
	Texto string `json:"text"`
}

type RequisicaoReacao struct {
	Reacao string `json:"reaction"`
}

type RequisicaoSalvarPost struct {
	Salvo bool `json:"saved"`
}

type PostFeedDetalhe struct {
	Identificador       string                   `json:"id"`
	AutorID             string                   `json:"author_id"`
	Autor               comum.PerfilPublicoAutor `json:"author"`
	Texto               string                   `json:"text"`
	Anexos              []AnexoPost              `json:"attachments"`
	TipoConteudo        string                   `json:"content_kind,omitempty"`
	EscopoPublicacao    string                   `json:"publish_scope,omitempty"`
	IDGrupoPublicacao   string                   `json:"publish_group_id,omitempty"`
	GosteiTotal         int                      `json:"likes_count"`
	DesgosteiTotal      int                      `json:"dislikes_count"`
	Comentarios         []ComentarioPost         `json:"comments"`
	MeuVoto             string                   `json:"my_reaction,omitempty"`
	Salvo               bool                     `json:"saved"`
	LinkCompartilhar    string                   `json:"share_link"`
	CriadoEm            string                   `json:"created_at"`
}

type FiltroListarPosts struct {
	Pagina            int
	Limite            int
	AutorID           string
	GruposDoUsuario   []string
	IncluirComentarios bool
}

type RespostaListaPosts struct {
	Itens   []PostFeedDetalhe `json:"items"`
	Total   int               `json:"total"`
	Pagina  int               `json:"page"`
	Limite  int               `json:"limit"`
	TemMais bool              `json:"has_more"`
}
