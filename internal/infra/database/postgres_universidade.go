package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	universidadeService "campus_connect_api/internal/modulos/universidade/structs"
)

func (p *Postgres) ListarAvisosUniversidade(ctx context.Context) ([]universidadeService.AvisoUniversidade, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []universidadeService.AvisoUniversidade
	for rows.Next() {
		var a universidadeService.AvisoUniversidade
		if err := rows.Scan(&a.Identificador, &a.Titulo, &a.Descricao); err != nil { return nil, err }
		out = append(out, a)
	}
	return out, rows.Err()
}
func (p *Postgres) ObterAvisoUniversidade(ctx context.Context, id string) (universidadeService.AvisoUniversidade, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade WHERE id=$1::uuid`
	var a universidadeService.AvisoUniversidade
	err := p.pool.QueryRow(ctx, sql, id).Scan(&a.Identificador, &a.Titulo, &a.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return universidadeService.AvisoUniversidade{}, false, nil }
		return universidadeService.AvisoUniversidade{}, false, err
	}
	return a, true, nil
}
func (p *Postgres) InserirAvisoUniversidade(ctx context.Context, criadoPor string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	tx, err := p.pool.Begin(ctx); if err != nil { return universidadeService.AvisoUniversidade{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO avisos_universidade (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil { return universidadeService.AvisoUniversidade{}, err }
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "notice", cartao, corpo.Titulo, "Aviso", corpo.Descricao, "Tipo", "universidade", id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil { return universidadeService.AvisoUniversidade{}, err }
	if err := tx.Commit(ctx); err != nil { return universidadeService.AvisoUniversidade{}, err }
	a, ok, err := p.ObterAvisoUniversidade(ctx, id); if err != nil || !ok { return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso") }
	return a, nil
}
func (p *Postgres) atualizarAvisoUniversidadeComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	if err := p.garantirDonoTabela(ctx, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil { return universidadeService.AvisoUniversidade{}, err }
	ct, err := p.pool.Exec(ctx, `UPDATE avisos_universidade SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil { return universidadeService.AvisoUniversidade{}, err }
	if ct.RowsAffected() == 0 { return universidadeService.AvisoUniversidade{}, ErrNaoEncontrado }
	a, ok, err := p.ObterAvisoUniversidade(ctx, id); if err != nil || !ok { return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso") }
	return a, nil
}
func (p *Postgres) removerAvisoUniversidadeComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil { return err }
	tx, err := p.pool.Begin(ctx); if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='notice' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM avisos_universidade WHERE id=$1::uuid`, id)
	if err != nil { return err }
	if ct.RowsAffected() == 0 { return ErrNaoEncontrado }
	return tx.Commit(ctx)
}
func (p *Postgres) AtualizarAvisoUniversidade(ctx context.Context, id, usuarioID string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return p.atualizarAvisoUniversidadeComPerfil(ctx, id, usuarioID, "padrao", corpo)
}
func (p *Postgres) AtualizarAvisoUniversidadeComoAdmin(ctx context.Context, id string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return p.atualizarAvisoUniversidadeComPerfil(ctx, id, "", "sistema_admin", corpo)
}
func (p *Postgres) RemoverAvisoUniversidade(ctx context.Context, id, usuarioID string) error {
	return p.removerAvisoUniversidadeComPerfil(ctx, id, usuarioID, "padrao")
}
func (p *Postgres) RemoverAvisoUniversidadeComoAdmin(ctx context.Context, id string) error {
	return p.removerAvisoUniversidadeComPerfil(ctx, id, "", "sistema_admin")
}
