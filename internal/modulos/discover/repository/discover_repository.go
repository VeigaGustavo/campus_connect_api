package repository

import (
	"context"
	"strings"

	discoverService "campus_connect_api/internal/modulos/discover/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type discoverRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoDiscoverRepository(pool *pgxpool.Pool) discoverService.DiscoverRepository {
	return &discoverRepositoryPostgres{pool: pool}
}

func (repositorio *discoverRepositoryPostgres) FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]discoverService.ItemDescobrir, error) {
	f := filtro
	if f == "" {
		f = "all"
	}
	var sql string
	var rows pgx.Rows
	var err error
	switch f {
	case "all":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes ORDER BY criado_em DESC`
	case "internships":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='internship' ORDER BY criado_em DESC`
	case "events":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='event' ORDER BY criado_em DESC`
	case "groups":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='study_group' ORDER BY criado_em DESC`
	case "projects":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='project' ORDER BY criado_em DESC`
	case "readings":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='reading' ORDER BY criado_em DESC`
	case "notices":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='notice' ORDER BY criado_em DESC`
	default:
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes ORDER BY criado_em DESC`
	}
	if len(gruposDoUsuario) == 0 {
		sql = strings.Replace(sql, " FROM feed_cartoes", " FROM feed_cartoes WHERE visibility_scope='all'", 1)
		rows, err = repositorio.pool.Query(contexto, sql)
	} else {
		sql = strings.Replace(sql, " FROM feed_cartoes", " FROM feed_cartoes WHERE (visibility_scope='all' OR (visibility_scope='group' AND visibility_group_id = ANY($1::text[])))", 1)
		rows, err = repositorio.pool.Query(contexto, sql, gruposDoUsuario)
	}
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
