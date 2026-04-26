package repository

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
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
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label, criado_por::text FROM leitura_semanal ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []leituraService.ItemLeituraSemanal
	for rows.Next() {
		var it leituraService.ItemLeituraSemanal
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta, &it.AutorID); err != nil {
			return nil, err
		}
		if err := carregarAutorPublico(contexto, repositorio.pool, it.AutorID, &it.Autor); err != nil {
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
	if err := repositoryutil.InserirCartaoFeedTx(contexto, tx, comum.FeedKindLeitura, "dsc-"+id, corpo.Titulo, corpo.Fonte, corpo.Resumo, "Leitura", corpo.Tipo, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
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
	return repositorio.atualizarLeituraSemanalComPerfil(contexto, id, usuarioID, comum.PerfilPadrao, corpo)
}

func (repositorio *leituraRepositoryPostgres) AtualizarLeituraSemanalComoAdmin(contexto context.Context, id string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return repositorio.atualizarLeituraSemanalComPerfil(contexto, id, "", comum.PerfilSistemaAdmin, corpo)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanal(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerLeituraSemanalComPerfil(contexto, id, usuarioID, comum.PerfilPadrao)
}

func (repositorio *leituraRepositoryPostgres) RemoverLeituraSemanalComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerLeituraSemanalComPerfil(contexto, id, "", comum.PerfilSistemaAdmin)
}

func (repositorio *leituraRepositoryPostgres) obterLeituraSemanal(contexto context.Context, id string) (leituraService.ItemLeituraSemanal, bool, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label, criado_por::text FROM leitura_semanal WHERE id=$1::uuid`
	var it leituraService.ItemLeituraSemanal
	var k string
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta, &it.AutorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return leituraService.ItemLeituraSemanal{}, false, nil
		}
		return leituraService.ItemLeituraSemanal{}, false, err
	}
	if err := carregarAutorPublico(contexto, repositorio.pool, it.AutorID, &it.Autor); err != nil {
		return leituraService.ItemLeituraSemanal{}, false, err
	}
	it.Tipo = leituraService.TipoItemLeitura(k)
	return it, true, nil
}

func (repositorio *leituraRepositoryPostgres) atualizarLeituraSemanalComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM leitura_semanal WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
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
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM leitura_semanal WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
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

func carregarAutorPublico(contexto context.Context, pool *pgxpool.Pool, autorID string, destino *comum.PerfilPublicoAutor) error {
	const sql = `SELECT u.id::text, u.nome, coalesce(u.avatar_image_url,''), pf.codigo
FROM usuarios u
JOIN perfis_usuario pf ON pf.id = u.perfil_id
WHERE u.id=$1::uuid`
	return pool.QueryRow(contexto, sql, autorID).Scan(&destino.Identificador, &destino.Nome, &destino.URLAvatar, &destino.Perfil)
}
