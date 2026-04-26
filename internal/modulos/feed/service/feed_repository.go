package service

import "context"

type FeedRepository interface {
	Feed(contexto context.Context, filtro string, gruposDoUsuario []string) ([]ItemFeed, error)
	CriarPost(contexto context.Context, criadoPor string, corpo RequisicaoCriarPost) (PostFeedDetalhe, error)
	ObterPost(contexto context.Context, postID, usuarioID string) (PostFeedDetalhe, bool, error)
	ListarComentariosPost(contexto context.Context, postID string) ([]ComentarioPost, error)
	CriarComentarioPost(contexto context.Context, postID, autorID string, corpo RequisicaoCriarComentario) (ComentarioPost, error)
	ReagirPost(contexto context.Context, postID, usuarioID, reacao string) error
	ReagirComentario(contexto context.Context, comentarioID, usuarioID, reacao string) error
	SalvarPost(contexto context.Context, postID, usuarioID string, salvo bool) error
}
