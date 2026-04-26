package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	projetoService "campus_connect_api/internal/modulos/projeto/structs"
)

func (p *Postgres) ListarProjetos(ctx context.Context) ([]projetoService.Projeto, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []projetoService.Projeto
	for rows.Next() {
		var pr projetoService.Projeto
		if err := rows.Scan(&pr.Identificador, &pr.Titulo, &pr.Descricao); err != nil { return nil, err }
		out = append(out, pr)
	}
	return out, rows.Err()
}
func (p *Postgres) ObterProjeto(ctx context.Context, id string) (projetoService.Projeto, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos WHERE id=$1::uuid`
	var pr projetoService.Projeto
	err := p.pool.QueryRow(ctx, sql, id).Scan(&pr.Identificador, &pr.Titulo, &pr.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return projetoService.Projeto{}, false, nil }
		return projetoService.Projeto{}, false, err
	}
	return pr, true, nil
}
func (p *Postgres) InserirProjeto(ctx context.Context, criadoPor string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	tx, err := p.pool.Begin(ctx); if err != nil { return projetoService.Projeto{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO projetos (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil { return projetoService.Projeto{}, err }
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "project", cartao, corpo.Titulo, "Projeto", corpo.Descricao, "Projeto", id, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil { return projetoService.Projeto{}, err }
	if err := tx.Commit(ctx); err != nil { return projetoService.Projeto{}, err }
	pr, ok, err := p.ObterProjeto(ctx, id); if err != nil || !ok { return projetoService.Projeto{}, errors.New("falha ao recarregar projeto") }
	return pr, nil
}
func (p *Postgres) atualizarProjetoComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	if err := p.garantirDonoTabela(ctx, "projetos", id, usuarioID, perfilCodigo); err != nil { return projetoService.Projeto{}, err }
	tx, err := p.pool.Begin(ctx); if err != nil { return projetoService.Projeto{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	ct, err := tx.Exec(ctx, `UPDATE projetos SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil { return projetoService.Projeto{}, err }
	if ct.RowsAffected() == 0 { return projetoService.Projeto{}, ErrNaoEncontrado }
	_, _ = tx.Exec(ctx, `UPDATE feed_cartoes SET titulo=$2, excerpt=$3 WHERE kind='project' AND reference_id=$1`, id, corpo.Titulo, corpo.Descricao)
	if err := tx.Commit(ctx); err != nil { return projetoService.Projeto{}, err }
	pr, ok, err := p.ObterProjeto(ctx, id); if err != nil || !ok { return projetoService.Projeto{}, errors.New("falha ao recarregar projeto") }
	return pr, nil
}
func (p *Postgres) removerProjetoComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "projetos", id, usuarioID, perfilCodigo); err != nil { return err }
	tx, err := p.pool.Begin(ctx); if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='project' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM projetos WHERE id=$1::uuid`, id)
	if err != nil { return err }
	if ct.RowsAffected() == 0 { return ErrNaoEncontrado }
	return tx.Commit(ctx)
}
func (p *Postgres) AtualizarProjeto(ctx context.Context, id, usuarioID string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return p.atualizarProjetoComPerfil(ctx, id, usuarioID, "padrao", corpo)
}
func (p *Postgres) AtualizarProjetoComoAdmin(ctx context.Context, id string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return p.atualizarProjetoComPerfil(ctx, id, "", "sistema_admin", corpo)
}
func (p *Postgres) RemoverProjeto(ctx context.Context, id, usuarioID string) error {
	return p.removerProjetoComPerfil(ctx, id, usuarioID, "padrao")
}
func (p *Postgres) RemoverProjetoComoAdmin(ctx context.Context, id string) error {
	return p.removerProjetoComPerfil(ctx, id, "", "sistema_admin")
}
