package repository

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	leituraService "campus_connect_api/internal/modulos/leitura/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type leituraRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoLeituraRepository(pool *pgxpool.Pool) leituraService.LeituraRepository {
	return &leituraRepositoryPostgres{pool: pool}
}

func (repositorio *leituraRepositoryPostgres) ListarLeituraSemanal(contexto context.Context) ([]leituraService.ItemLeituraSemanal, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []leituraService.ItemLeituraSemanal
	for rows.Next() {
		var it leituraService.ItemLeituraSemanal
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta); err != nil {
			return nil, err
		}
		it.Tipo = leituraService.TipoItemLeitura(k)
		out = append(out, it)
	}
	return out, rows.Err()
}

func (repositorio *leituraRepositoryPostgres) InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const ins = `INSERT INTO leitura_semanal (kind, titulo, source, excerpt, image_url, meta_label, criado_por) VALUES ($1::varchar,$2,$3,$4,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(contexto, ins, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta, criadoPor).Scan(&id); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	if err := inserirCartaoFeedTx(contexto, tx, "reading", "dsc-"+id, corpo.Titulo, corpo.Fonte, corpo.Resumo, "Leitura", corpo.Tipo, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	it, ok, err := repositorio.obterLeituraSemanal(contexto, id)
	if err != nil || !ok {
		return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura")
	}
	return it, nil
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanal(contexto context.Context, id, usuarioID string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.atualizarLeituraSemanalComPerfil(contexto, id, usuarioID, "padrao", corpo)
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanalComoAdmin(contexto context.Context, id string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.atualizarLeituraSemanalComPerfil(contexto, id, "", "sistema_admin", corpo)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanal(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerLeituraSemanalComPerfil(contexto, id, usuarioID, "padrao")
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanalComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerLeituraSemanalComPerfil(contexto, id, "", "sistema_admin")
}

func (repositorio *leituraRepositoryPostgres) obterLeituraSemanal(contexto context.Context, id string) (leituraService.ItemLeituraSemanal, bool, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal WHERE id=$1::uuid`
	var it leituraService.ItemLeituraSemanal
	var k string
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return leituraService.ItemLeituraSemanal{}, false, nil
		}
		return leituraService.ItemLeituraSemanal{}, false, err
	}
	it.Tipo = leituraService.TipoItemLeitura(k)
	return it, true, nil
}

func (repositorio *leituraRepositoryPostgres) atualizarLeituraSemanalComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	if err := garantirDonoTabela(contexto, repositorio.pool, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	ct, err := repositorio.pool.Exec(contexto, `UPDATE leitura_semanal SET kind=$2::varchar, titulo=$3, source=$4, excerpt=$5, image_url=$6, meta_label=$7, atualizado_em=now() WHERE id=$1::uuid`,
		id, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta)
	if err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	if ct.RowsAffected() == 0 {
		return leituraService.ItemLeituraSemanal{}, comum.ErrNaoEncontrado
	}
	it, ok, err := repositorio.obterLeituraSemanal(contexto, id)
	if err != nil || !ok {
		return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura")
	}
	return it, nil
}

func (repositorio *leituraRepositoryPostgres) removerLeituraSemanalComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := garantirDonoTabela(contexto, repositorio.pool, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='reading' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM leitura_semanal WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return tx.Commit(contexto)
}

func garantirDonoTabela(ctx context.Context, pool *pgxpool.Pool, tabela, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return nil
	}
	var q string
	switch tabela {
	case "leitura_semanal":
		q = `SELECT criado_por::text FROM leitura_semanal WHERE id=$1::uuid`
	default:
		return comum.ErrProibido
	}
	var dono string
	err := pool.QueryRow(ctx, q, id).Scan(&dono)
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

func inserirCartaoFeedTx(ctx context.Context, tx pgx.Tx, kind, cartaoID, titulo, subtitulo, excerpt, metaPri, metaSec, ref, scope, groupID string) error {
	if scope != "group" {
		scope = "all"
		groupID = ""
	}
	const sql = `
INSERT INTO feed_cartoes (id, kind, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id, visibility_scope, visibility_group_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := tx.Exec(ctx, sql, cartaoID, kind, titulo, subtitulo, excerpt, metaPri, metaSec, ref, scope, nullSeVazio(groupID))
	return err
}

func nullSeVazio(s string) any {
	if s == "" {
		return nil
	}
	return s
}
