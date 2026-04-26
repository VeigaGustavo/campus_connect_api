package model

type ItemFeed struct {
	Identificador  string            `json:"id"`
	Categoria      CategoriaItemFeed `json:"kind"`
	Titulo         string            `json:"title"`
	Subtitulo      string            `json:"subtitle"`
	Resumo         string            `json:"excerpt"`
	MetaPrincipal  string            `json:"meta_primary"`
	MetaSecundaria string            `json:"meta_secondary"`
	IDReferencia   string            `json:"reference_id"`
}

type RespostaFeed struct {
	Itens []ItemFeed `json:"items"`
}

type AnexoPost struct {
	Tipo string `json:"type"`
	URL  string `json:"url"`
	Nome string `json:"name,omitempty"`
}

type PerfilAutor struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	URLAvatar     string `json:"avatar_image_url"`
	Perfil        string `json:"role"`
}

type RequisicaoCriarPost struct {
	Texto             string      `json:"text"`
	Anexos            []AnexoPost `json:"attachments"`
	EscopoPublicacao  string      `json:"publish_scope"`
	IDGrupoPublicacao string      `json:"publish_group_id"`
}

type ComentarioPost struct {
	Identificador  string      `json:"id"`
	PostID         string      `json:"post_id"`
	AutorID        string      `json:"author_id"`
	Autor          PerfilAutor `json:"author"`
	Texto          string      `json:"text"`
	GosteiTotal    int         `json:"likes_count"`
	DesgosteiTotal int         `json:"dislikes_count"`
	CriadoEm       string      `json:"created_at"`
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
	Identificador    string           `json:"id"`
	AutorID          string           `json:"author_id"`
	Autor            PerfilAutor      `json:"author"`
	Texto            string           `json:"text"`
	Anexos           []AnexoPost      `json:"attachments"`
	GosteiTotal      int              `json:"likes_count"`
	DesgosteiTotal   int              `json:"dislikes_count"`
	Comentarios      []ComentarioPost `json:"comments"`
	MeuVoto          string           `json:"my_reaction,omitempty"`
	Salvo            bool             `json:"saved"`
	LinkCompartilhar string           `json:"share_link"`
	CriadoEm         string           `json:"created_at"`
}
