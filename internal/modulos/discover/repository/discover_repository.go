package repository

import (
	"context"
	"fmt"

	comum "campus_connect_api/internal/modulos/comum"
	discoverService "campus_connect_api/internal/modulos/discover/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type discoverRepositoryPostgres struct {
	pool *pgxpool.Pool
}

const selectFeedDescobrir = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes`

var feedKindPorFiltro = map[string]string{
	comum.FiltroDescobrirEstagios: comum.FeedKindEstagio,
	comum.FiltroDescobrirEventos:  comum.FeedKindEvento,
	comum.FiltroDescobrirGrupos:   comum.FeedKindGrupoEstudo,
	comum.FiltroDescobrirProjetos: comum.FeedKindProjeto,
	comum.FiltroDescobrirLeituras: comum.FeedKindLeitura,
	comum.FiltroDescobrirAvisos:   comum.FeedKindAviso,
}

func NovoDiscoverRepository(pool *pgxpool.Pool) discoverService.DiscoverRepository {
	return &discoverRepositoryPostgres{pool: pool}
}

func (repositorio *discoverRepositoryPostgres) FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]discoverService.ItemDescobrir, error) {
	f := filtro
	if f == "" {
		f = comum.FiltroDescobrirTodos
	}
	sql := selectFeedDescobrir + " WHERE 1=1"
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
	var out []discoverService.ItemDescobrir
	for rows.Next() {
		var it discoverService.ItemDescobrir
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Subtitulo, &it.Resumo, &it.MetaPrincipal, &it.MetaSecundaria, &it.IDReferencia); err != nil {
			return nil, err
		}
		it.Categoria = discoverService.CategoriaItemDescobrir(k)
		out = append(out, it)
	}
	return out, rows.Err()
}
