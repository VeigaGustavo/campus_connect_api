package repositoryutil

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GarantirDonoOuAdmin(ctx context.Context, pool *pgxpool.Pool, ownerQuery, id, usuarioID, perfil string) error {
	if perfil == comum.PerfilSistemaAdmin {
		return nil
	}

	var dono string
	err := pool.QueryRow(ctx, ownerQuery, id).Scan(&dono)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return comum.ErrNaoEncontrado
		}
		return err
	}
	if dono != usuarioID {
		return comum.ErrProibido
	}
	return nil
}

func InserirCartaoFeedTx(ctx context.Context, tx pgx.Tx, kind, cartaoID, titulo, subtitulo, excerpt, metaPri, metaSec, ref, scope, groupID string) error {
	if scope != comum.VisibilidadeGrupo {
		scope = comum.VisibilidadeTodos
		groupID = ""
	}
	const sql = `
INSERT INTO feed_cartoes (id, kind, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id, visibility_scope, visibility_group_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := tx.Exec(ctx, sql, cartaoID, kind, titulo, subtitulo, excerpt, metaPri, metaSec, ref, scope, NullSeVazio(groupID))
	return err
}

func NullSeVazio(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// CarregarPerfilPublicoAutor preenche nome, avatar e codigo de perfil do usuario (uso em listagens de conteudo).
func CarregarPerfilPublicoAutor(ctx context.Context, pool *pgxpool.Pool, usuarioID string, destino *comum.PerfilPublicoAutor) error {
	const sql = `SELECT u.id::text, u.nome, coalesce(u.avatar_image_url,''), pf.codigo
FROM usuarios u
JOIN perfis_usuario pf ON pf.id = u.perfil_id
WHERE u.id=$1::uuid`
	return pool.QueryRow(ctx, sql, usuarioID).Scan(&destino.Identificador, &destino.Nome, &destino.URLAvatar, &destino.Perfil)
}
