package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
	feedService "campus_connect_api/internal/modulos/feed/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type feedRepositoryPostgres struct {
	pool *pgxpool.Pool
}

const selectFeed = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes`

var feedKindPorFiltro = map[string]string{
	comum.FiltroDescobrirEstagios: comum.FeedKindEstagio,
	comum.FiltroDescobrirEventos:  comum.FeedKindEvento,
	comum.FiltroDescobrirGrupos:   comum.FeedKindGrupoEstudo,
	comum.FiltroDescobrirProjetos: comum.FeedKindProjeto,
	comum.FiltroDescobrirLeituras: comum.FeedKindLeitura,
	comum.FiltroDescobrirAvisos:   comum.FeedKindAviso,
	"posts":                       "post",
}

func NovoFeedRepository(pool *pgxpool.Pool) feedService.FeedRepository {
	return &feedRepositoryPostgres{pool: pool}
}

func (repositorio *feedRepositoryPostgres) Feed(contexto context.Context, filtro string, gruposDoUsuario []string) ([]feedService.ItemFeed, error) {
	f := filtro
	if f == "" {
		f = comum.FiltroDescobrirTodos
	}
	sql := selectFeed + " WHERE 1=1"
	args := make([]any, 0, 4)
	var rows pgx.Rows
	var err error

	if kind, ok := feedKindPorFiltro[f]; ok {
		kindPos := len(args) + 1
		args = append(args, kind)
		sql += fmt.Sprintf(" AND kind=$%d", kindPos)
	}

	if len(gruposDoUsuario) == 0 {
		visPos := len(args) + 1
		args = append(args, comum.VisibilidadeTodos)
		sql += fmt.Sprintf(" AND visibility_scope=$%d", visPos)
	} else {
		todosPos := len(args) + 1
		grupoPos := len(args) + 2
		gruposPos := len(args) + 3
		args = append(args, comum.VisibilidadeTodos, comum.VisibilidadeGrupo, gruposDoUsuario)
		sql += fmt.Sprintf(" AND (visibility_scope=$%d OR (visibility_scope=$%d AND visibility_group_id = ANY($%d::text[])))", todosPos, grupoPos, gruposPos)
	}

	sql += " ORDER BY criado_em DESC"
	rows, err = repositorio.pool.Query(contexto, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []feedService.ItemFeed
	for rows.Next() {
		var it feedService.ItemFeed
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Subtitulo, &it.Resumo, &it.MetaPrincipal, &it.MetaSecundaria, &it.IDReferencia); err != nil {
			return nil, err
		}
		it.Categoria = feedService.CategoriaItemFeed(k)
		out = append(out, it)
	}

	return out, rows.Err()
}

func (repositorio *feedRepositoryPostgres) CriarPost(contexto context.Context, criadoPor string, corpo feedService.RequisicaoCriarPost) (feedService.PostFeedDetalhe, error) {
	anexosJSON, err := json.Marshal(corpo.Anexos)
	if err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()

	const insPost = `INSERT INTO feed_posts (author_id, body_text, attachments, share_link) VALUES ($1::uuid,$2,$3::jsonb,'') RETURNING id::text, criado_em`
	var postID string
	var criadoEm time.Time
	if err := tx.QueryRow(contexto, insPost, criadoPor, corpo.Texto, anexosJSON).Scan(&postID, &criadoEm); err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	shareLink := "/posts/" + postID
	if _, err := tx.Exec(contexto, `UPDATE feed_posts SET share_link=$2 WHERE id=$1::uuid`, postID, shareLink); err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	scope := corpo.EscopoPublicacao
	if scope != comum.VisibilidadeGrupo {
		scope = comum.VisibilidadeTodos
		corpo.IDGrupoPublicacao = ""
	}
	const insCard = `INSERT INTO feed_cartoes (id, kind, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id, visibility_scope, visibility_group_id)
VALUES ($1,'post',$2,$3,$4,$5,$6,$7,$8,$9)`
	if _, err := tx.Exec(contexto, insCard, "dsc-"+postID, "Novo post", "Comunidade", corpo.Texto, "Post", "social", postID, scope, nullSeVazio(corpo.IDGrupoPublicacao)); err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return feedService.PostFeedDetalhe{}, err
	}
	post, _, err := repositorio.ObterPost(contexto, postID, criadoPor)
	return post, err
}

func (repositorio *feedRepositoryPostgres) ObterPost(contexto context.Context, postID, usuarioID string) (feedService.PostFeedDetalhe, bool, error) {
	const sql = `SELECT id::text, author_id::text, body_text, attachments, share_link, criado_em FROM feed_posts WHERE id=$1::uuid`
	var post feedService.PostFeedDetalhe
	var anexosJSON []byte
	var criadoEm time.Time
	err := repositorio.pool.QueryRow(contexto, sql, postID).Scan(&post.Identificador, &post.AutorID, &post.Texto, &anexosJSON, &post.LinkCompartilhar, &criadoEm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return feedService.PostFeedDetalhe{}, false, nil
		}
		return feedService.PostFeedDetalhe{}, false, err
	}
	_ = json.Unmarshal(anexosJSON, &post.Anexos)
	post.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
	if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, post.AutorID, &post.Autor); err != nil {
		return feedService.PostFeedDetalhe{}, false, err
	}
	post.Comentarios, err = repositorio.ListarComentariosPost(contexto, postID)
	if err != nil {
		return feedService.PostFeedDetalhe{}, false, err
	}
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM feed_post_reacoes WHERE post_id=$1::uuid AND reaction='like'`, postID).Scan(&post.GosteiTotal); err != nil {
		return feedService.PostFeedDetalhe{}, false, err
	}
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM feed_post_reacoes WHERE post_id=$1::uuid AND reaction='dislike'`, postID).Scan(&post.DesgosteiTotal); err != nil {
		return feedService.PostFeedDetalhe{}, false, err
	}
	if usuarioID != "" {
		_ = repositorio.pool.QueryRow(contexto, `SELECT reaction FROM feed_post_reacoes WHERE post_id=$1::uuid AND user_id=$2::uuid`, postID, usuarioID).Scan(&post.MeuVoto)
		var total int
		if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM feed_posts_salvos WHERE post_id=$1::uuid AND user_id=$2::uuid`, postID, usuarioID).Scan(&total); err == nil {
			post.Salvo = total > 0
		}
	}
	return post, true, nil
}

func (repositorio *feedRepositoryPostgres) ListarComentariosPost(contexto context.Context, postID string) ([]feedService.ComentarioPost, error) {
	const sql = `SELECT c.id::text, c.post_id::text, c.author_id::text, c.body_text, c.criado_em,
COALESCE(sum(CASE WHEN r.reaction='like' THEN 1 ELSE 0 END),0) AS likes_count,
COALESCE(sum(CASE WHEN r.reaction='dislike' THEN 1 ELSE 0 END),0) AS dislikes_count
FROM feed_comentarios c
LEFT JOIN feed_comentario_reacoes r ON r.comment_id = c.id
WHERE c.post_id=$1::uuid
GROUP BY c.id
ORDER BY c.criado_em ASC`
	rows, err := repositorio.pool.Query(contexto, sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []feedService.ComentarioPost
	for rows.Next() {
		var it feedService.ComentarioPost
		var criadoEm time.Time
		if err := rows.Scan(&it.Identificador, &it.PostID, &it.AutorID, &it.Texto, &criadoEm, &it.GosteiTotal, &it.DesgosteiTotal); err != nil {
			return nil, err
		}
		if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, it.AutorID, &it.Autor); err != nil {
			return nil, err
		}
		it.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
		out = append(out, it)
	}
	return out, rows.Err()
}

func (repositorio *feedRepositoryPostgres) CriarComentarioPost(contexto context.Context, postID, autorID string, corpo feedService.RequisicaoCriarComentario) (feedService.ComentarioPost, error) {
	const sql = `INSERT INTO feed_comentarios (post_id, author_id, body_text) VALUES ($1::uuid,$2::uuid,$3) RETURNING id::text, criado_em`
	var out feedService.ComentarioPost
	var criadoEm time.Time
	if err := repositorio.pool.QueryRow(contexto, sql, postID, autorID, corpo.Texto).Scan(&out.Identificador, &criadoEm); err != nil {
		return feedService.ComentarioPost{}, err
	}
	out.PostID = postID
	out.AutorID = autorID
	if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, autorID, &out.Autor); err != nil {
		return feedService.ComentarioPost{}, err
	}
	out.Texto = corpo.Texto
	out.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
	return out, nil
}

func (repositorio *feedRepositoryPostgres) ReagirPost(contexto context.Context, postID, usuarioID, reacao string) error {
	if reacao == "" {
		_, err := repositorio.pool.Exec(contexto, `DELETE FROM feed_post_reacoes WHERE post_id=$1::uuid AND user_id=$2::uuid`, postID, usuarioID)
		return err
	}
	const sql = `INSERT INTO feed_post_reacoes (post_id, user_id, reaction) VALUES ($1::uuid,$2::uuid,$3)
ON CONFLICT (post_id, user_id) DO UPDATE SET reaction=excluded.reaction, atualizado_em=now()`
	_, err := repositorio.pool.Exec(contexto, sql, postID, usuarioID, reacao)
	return err
}

func (repositorio *feedRepositoryPostgres) ReagirComentario(contexto context.Context, comentarioID, usuarioID, reacao string) error {
	if reacao == "" {
		_, err := repositorio.pool.Exec(contexto, `DELETE FROM feed_comentario_reacoes WHERE comment_id=$1::uuid AND user_id=$2::uuid`, comentarioID, usuarioID)
		return err
	}
	const sql = `INSERT INTO feed_comentario_reacoes (comment_id, user_id, reaction) VALUES ($1::uuid,$2::uuid,$3)
ON CONFLICT (comment_id, user_id) DO UPDATE SET reaction=excluded.reaction, atualizado_em=now()`
	_, err := repositorio.pool.Exec(contexto, sql, comentarioID, usuarioID, reacao)
	return err
}

func (repositorio *feedRepositoryPostgres) SalvarPost(contexto context.Context, postID, usuarioID string, salvo bool) error {
	if !salvo {
		_, err := repositorio.pool.Exec(contexto, `DELETE FROM feed_posts_salvos WHERE post_id=$1::uuid AND user_id=$2::uuid`, postID, usuarioID)
		return err
	}
	const sql = `INSERT INTO feed_posts_salvos (post_id, user_id) VALUES ($1::uuid,$2::uuid) ON CONFLICT (post_id, user_id) DO NOTHING`
	_, err := repositorio.pool.Exec(contexto, sql, postID, usuarioID)
	return err
}

func nullSeVazio(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

