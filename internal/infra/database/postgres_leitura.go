package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	leituraService "campus_connect_api/internal/modulos/leitura/structs"
)

func (p *Postgres) ListarLeituraSemanal(ctx context.Context) ([]leituraService.ItemLeituraSemanal, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []leituraService.ItemLeituraSemanal
	for rows.Next() {
		var it leituraService.ItemLeituraSemanal
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta); err != nil { return nil, err }
		it.Tipo = leituraService.TipoItemLeitura(k)
		out = append(out, it)
	}
	return out, rows.Err()
}
func (p *Postgres) ObterLeituraSemanal(ctx context.Context, id string) (leituraService.ItemLeituraSemanal, bool, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal WHERE id=$1::uuid`
	var it leituraService.ItemLeituraSemanal
	var k string
	err := p.pool.QueryRow(ctx, sql, id).Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return leituraService.ItemLeituraSemanal{}, false, nil }
		return leituraService.ItemLeituraSemanal{}, false, err
	}
	it.Tipo = leituraService.TipoItemLeitura(k)
	return it, true, nil
}
func (p *Postgres) InserirLeituraSemanal(ctx context.Context, criadoPor string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	tx, err := p.pool.Begin(ctx); if err != nil { return leituraService.ItemLeituraSemanal{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO leitura_semanal (kind, titulo, source, excerpt, image_url, meta_label, criado_por) VALUES ($1::varchar,$2,$3,$4,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta, criadoPor).Scan(&id); err != nil { return leituraService.ItemLeituraSemanal{}, err }
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "reading", cartao, corpo.Titulo, corpo.Fonte, corpo.Resumo, "Leitura", corpo.Tipo, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil { return leituraService.ItemLeituraSemanal{}, err }
	if err := tx.Commit(ctx); err != nil { return leituraService.ItemLeituraSemanal{}, err }
	it, ok, err := p.ObterLeituraSemanal(ctx, id); if err != nil || !ok { return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura") }
	return it, nil
}
func (p *Postgres) atualizarLeituraSemanalComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	if err := p.garantirDonoTabela(ctx, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil { return leituraService.ItemLeituraSemanal{}, err }
	ct, err := p.pool.Exec(ctx, `UPDATE leitura_semanal SET kind=$2::varchar, titulo=$3, source=$4, excerpt=$5, image_url=$6, meta_label=$7, atualizado_em=now() WHERE id=$1::uuid`,
		id, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta)
	if err != nil { return leituraService.ItemLeituraSemanal{}, err }
	if ct.RowsAffected() == 0 { return leituraService.ItemLeituraSemanal{}, ErrNaoEncontrado }
	it, ok, err := p.ObterLeituraSemanal(ctx, id); if err != nil || !ok { return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura") }
	return it, nil
}
func (p *Postgres) removerLeituraSemanalComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil { return err }
	tx, err := p.pool.Begin(ctx); if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='reading' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM leitura_semanal WHERE id=$1::uuid`, id)
	if err != nil { return err }
	if ct.RowsAffected() == 0 { return ErrNaoEncontrado }
	return tx.Commit(ctx)
}
func (p *Postgres) AtualizarLeituraSemanal(ctx context.Context, id, usuarioID string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return p.atualizarLeituraSemanalComPerfil(ctx, id, usuarioID, "padrao", corpo)
}
func (p *Postgres) AtualizarLeituraSemanalComoAdmin(ctx context.Context, id string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	return p.atualizarLeituraSemanalComPerfil(ctx, id, "", "sistema_admin", corpo)
}
func (p *Postgres) RemoverLeituraSemanal(ctx context.Context, id, usuarioID string) error {
	return p.removerLeituraSemanalComPerfil(ctx, id, usuarioID, "padrao")
}
func (p *Postgres) RemoverLeituraSemanalComoAdmin(ctx context.Context, id string) error {
	return p.removerLeituraSemanalComPerfil(ctx, id, "", "sistema_admin")
}
