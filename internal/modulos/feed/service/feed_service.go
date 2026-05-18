package service

import "context"

type FeedService struct {
	repositorio FeedRepository
}

func NovoFeedService(repositorio FeedRepository) *FeedService {
	return &FeedService{repositorio: repositorio}
}

func (servico *FeedService) Feed(contexto context.Context, filtro string, gruposDoUsuario []string) (RespostaFeed, error) {
	itens, err := servico.repositorio.Feed(contexto, filtro, gruposDoUsuario)
	if err != nil {
		return RespostaFeed{}, err
	}
	return RespostaFeed{Itens: itens}, nil
}

func (servico *FeedService) ListarPosts(contexto context.Context, usuarioID string, filtro FiltroListarPosts) (RespostaListaPosts, error) {
	return servico.repositorio.ListarPosts(contexto, usuarioID, filtro)
}

func (servico *FeedService) CriarPost(contexto context.Context, criadoPor string, corpo RequisicaoCriarPost) (PostFeedDetalhe, error) {
	if err := validarCriarPost(corpo); err != nil {
		return PostFeedDetalhe{}, err
	}
	corpo.Anexos = normalizarAnexos(corpo.Anexos)
	if corpo.EscopoPublicacao == "" {
		corpo.EscopoPublicacao = "all"
	}
	return servico.repositorio.CriarPost(contexto, criadoPor, corpo)
}

func (servico *FeedService) ObterPost(contexto context.Context, postID, usuarioID string) (PostFeedDetalhe, bool, error) {
	return servico.repositorio.ObterPost(contexto, postID, usuarioID)
}

func (servico *FeedService) ListarComentariosPost(contexto context.Context, postID string) ([]ComentarioPost, error) {
	return servico.repositorio.ListarComentariosPost(contexto, postID)
}

func (servico *FeedService) CriarComentarioPost(contexto context.Context, postID, autorID string, corpo RequisicaoCriarComentario) (ComentarioPost, error) {
	return servico.repositorio.CriarComentarioPost(contexto, postID, autorID, corpo)
}

func (servico *FeedService) ReagirPost(contexto context.Context, postID, usuarioID, reacao string) error {
	return servico.repositorio.ReagirPost(contexto, postID, usuarioID, reacao)
}

func (servico *FeedService) ReagirComentario(contexto context.Context, comentarioID, usuarioID, reacao string) error {
	return servico.repositorio.ReagirComentario(contexto, comentarioID, usuarioID, reacao)
}

func (servico *FeedService) SalvarPost(contexto context.Context, postID, usuarioID string, salvo bool) error {
	return servico.repositorio.SalvarPost(contexto, postID, usuarioID, salvo)
}
