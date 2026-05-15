package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
	feedService "campus_connect_api/internal/modulos/feed/service"
)

func (repositorio *feedRepositoryPostgres) ListarPosts(contexto context.Context, usuarioID string, filtro feedService.FiltroListarPosts) (feedService.RespostaListaPosts, error) {
	pagina := filtro.Pagina
	if pagina < 1 {
		pagina = 1
	}
	limite := filtro.Limite
	if limite < 1 {
		limite = 20
	}
	if limite > 50 {
		limite = 50
	}
	offset := (pagina - 1) * limite

	baseFrom := `FROM feed_posts p
INNER JOIN feed_cartoes c ON c.kind = 'post' AND c.reference_id = p.id::text`
	where := []string{"1=1"}
	args := make([]any, 0, 8)

	if autor := strings.TrimSpace(filtro.AutorID); autor != "" {
		args = append(args, autor)
		where = append(where, fmt.Sprintf("p.author_id = $%d::uuid", len(args)))
	}

	if len(filtro.GruposDoUsuario) == 0 {
		args = append(args, comum.VisibilidadeTodos)
		where = append(where, fmt.Sprintf("c.visibility_scope = $%d", len(args)))
	} else {
		args = append(args, comum.VisibilidadeTodos, comum.VisibilidadeGrupo, filtro.GruposDoUsuario)
		todosPos := len(args) - 2
		grupoPos := len(args) - 1
		gruposPos := len(args)
		where = append(where, fmt.Sprintf("(c.visibility_scope = $%d OR (c.visibility_scope = $%d AND c.visibility_group_id = ANY($%d::text[])))", todosPos, grupoPos, gruposPos))
	}

	whereSQL := strings.Join(where, " AND ")

	var total int
	if err := repositorio.pool.QueryRow(contexto, "SELECT count(1) "+baseFrom+" WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return feedService.RespostaListaPosts{}, err
	}

	argsPag := append(append([]any{}, args...), limite, offset)
	limPos := len(argsPag) - 1
	offPos := len(argsPag)
	sql := `SELECT p.id::text, p.author_id::text, p.body_text, p.attachments, p.share_link, p.criado_em, coalesce(p.content_kind,''),
coalesce(c.visibility_scope,'all'), coalesce(c.visibility_group_id::text,'') ` + baseFrom + ` WHERE ` + whereSQL +
		fmt.Sprintf(" ORDER BY p.criado_em DESC LIMIT $%d OFFSET $%d", limPos, offPos)

	rows, err := repositorio.pool.Query(contexto, sql, argsPag...)
	if err != nil {
		return feedService.RespostaListaPosts{}, err
	}
	defer rows.Close()

	posts := make([]feedService.PostFeedDetalhe, 0, limite)
	ids := make([]string, 0, limite)
	for rows.Next() {
		var post feedService.PostFeedDetalhe
		var anexosJSON []byte
		var criadoEm time.Time
		if err := rows.Scan(&post.Identificador, &post.AutorID, &post.Texto, &anexosJSON, &post.LinkCompartilhar, &criadoEm, &post.TipoConteudo, &post.EscopoPublicacao, &post.IDGrupoPublicacao); err != nil {
			return feedService.RespostaListaPosts{}, err
		}
		_ = json.Unmarshal(anexosJSON, &post.Anexos)
		if post.Anexos == nil {
			post.Anexos = []feedService.AnexoPost{}
		}
		post.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
		post.Comentarios = []feedService.ComentarioPost{}
		posts = append(posts, post)
		ids = append(ids, post.Identificador)
	}
	if err := rows.Err(); err != nil {
		return feedService.RespostaListaPosts{}, err
	}

	if err := repositorio.enriquecerPostsLista(contexto, usuarioID, posts, ids, filtro.IncluirComentarios); err != nil {
		return feedService.RespostaListaPosts{}, err
	}

	return feedService.RespostaListaPosts{
		Itens:   posts,
		Total:   total,
		Pagina:  pagina,
		Limite:  limite,
		TemMais: offset+len(posts) < total,
	}, nil
}

func (repositorio *feedRepositoryPostgres) enriquecerPostsLista(contexto context.Context, usuarioID string, posts []feedService.PostFeedDetalhe, ids []string, incluirComentarios bool) error {
	if len(posts) == 0 {
		return nil
	}

	autores := make(map[string]struct{})
	for _, p := range posts {
		autores[p.AutorID] = struct{}{}
	}
	cacheAutor := make(map[string]comum.PerfilPublicoAutor, len(autores))
	for autorID := range autores {
		var autor comum.PerfilPublicoAutor
		if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, autorID, &autor); err != nil {
			return err
		}
		cacheAutor[autorID] = autor
	}

	likes, dislikes, err := repositorio.contagensReacoesPosts(contexto, ids)
	if err != nil {
		return err
	}

	minhas, err := repositorio.reacoesUsuarioPosts(contexto, usuarioID, ids)
	if err != nil {
		return err
	}
	salvos, err := repositorio.postsSalvosUsuario(contexto, usuarioID, ids)
	if err != nil {
		return err
	}

	var comentarios map[string][]feedService.ComentarioPost
	if incluirComentarios {
		comentarios = make(map[string][]feedService.ComentarioPost, len(ids))
		for _, id := range ids {
			lista, err := repositorio.ListarComentariosPost(contexto, id)
			if err != nil {
				return err
			}
			if lista == nil {
				lista = []feedService.ComentarioPost{}
			}
			comentarios[id] = lista
		}
	}

	for i := range posts {
		posts[i].Autor = cacheAutor[posts[i].AutorID]
		posts[i].GosteiTotal = likes[posts[i].Identificador]
		posts[i].DesgosteiTotal = dislikes[posts[i].Identificador]
		posts[i].MeuVoto = minhas[posts[i].Identificador]
		posts[i].Salvo = salvos[posts[i].Identificador]
		if incluirComentarios {
			posts[i].Comentarios = comentarios[posts[i].Identificador]
		}
	}
	return nil
}

func (repositorio *feedRepositoryPostgres) contagensReacoesPosts(contexto context.Context, ids []string) (likes, dislikes map[string]int, err error) {
	likes = make(map[string]int, len(ids))
	dislikes = make(map[string]int, len(ids))
	const sql = `SELECT post_id::text, reaction, count(1) FROM feed_post_reacoes WHERE post_id = ANY($1::uuid[]) GROUP BY post_id, reaction`
	rows, err := repositorio.pool.Query(contexto, sql, ids)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID, reaction string
		var total int
		if err := rows.Scan(&postID, &reaction, &total); err != nil {
			return nil, nil, err
		}
		switch reaction {
		case "like":
			likes[postID] = total
		case "dislike":
			dislikes[postID] = total
		}
	}
	return likes, dislikes, rows.Err()
}

func (repositorio *feedRepositoryPostgres) reacoesUsuarioPosts(contexto context.Context, usuarioID string, ids []string) (map[string]string, error) {
	out := make(map[string]string)
	if usuarioID == "" {
		return out, nil
	}
	const sql = `SELECT post_id::text, reaction FROM feed_post_reacoes WHERE user_id = $1::uuid AND post_id = ANY($2::uuid[])`
	rows, err := repositorio.pool.Query(contexto, sql, usuarioID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID, reaction string
		if err := rows.Scan(&postID, &reaction); err != nil {
			return nil, err
		}
		out[postID] = reaction
	}
	return out, rows.Err()
}

func (repositorio *feedRepositoryPostgres) postsSalvosUsuario(contexto context.Context, usuarioID string, ids []string) (map[string]bool, error) {
	out := make(map[string]bool)
	if usuarioID == "" {
		return out, nil
	}
	const sql = `SELECT post_id::text FROM feed_posts_salvos WHERE user_id = $1::uuid AND post_id = ANY($2::uuid[])`
	rows, err := repositorio.pool.Query(contexto, sql, usuarioID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		out[postID] = true
	}
	return out, rows.Err()
}

func (repositorio *feedRepositoryPostgres) carregarEscopoPost(contexto context.Context, postID string, post *feedService.PostFeedDetalhe) {
	_ = repositorio.pool.QueryRow(contexto, `
SELECT coalesce(visibility_scope,'all'), coalesce(visibility_group_id::text,'')
FROM feed_cartoes WHERE kind = 'post' AND reference_id = $1`, postID).Scan(&post.EscopoPublicacao, &post.IDGrupoPublicacao)
}
