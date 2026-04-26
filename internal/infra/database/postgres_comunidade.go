package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	comunidadeService "campus_connect_api/internal/modulos/comunidade/structs"
)

func (p *Postgres) ListarComunidades(ctx context.Context) ([]comunidadeService.Comunidade, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []comunidadeService.Comunidade
	for rows.Next() {
		var c comunidadeService.Comunidade
		if err := rows.Scan(&c.Identificador, &c.Nome, &c.Tipo, &c.Descricao); err != nil { return nil, err }
		out = append(out, c)
	}
	return out, rows.Err()
}
func (p *Postgres) ObterComunidade(ctx context.Context, id string) (comunidadeService.Comunidade, bool, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades WHERE id=$1::uuid`
	var c comunidadeService.Comunidade
	err := p.pool.QueryRow(ctx, sql, id).Scan(&c.Identificador, &c.Nome, &c.Tipo, &c.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return comunidadeService.Comunidade{}, false, nil }
		return comunidadeService.Comunidade{}, false, err
	}
	return c, true, nil
}
func (p *Postgres) InserirComunidade(ctx context.Context, criadoPor string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	const ins = `INSERT INTO comunidades (nome, kind, description, criado_por) VALUES ($1,$2,$3,$4::uuid) RETURNING id::text`
	var id string
	if err := p.pool.QueryRow(ctx, ins, corpo.Nome, corpo.Tipo, corpo.Descricao, criadoPor).Scan(&id); err != nil { return comunidadeService.Comunidade{}, err }
	c, ok, err := p.ObterComunidade(ctx, id); if err != nil || !ok { return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade") }
	return c, nil
}
func (p *Postgres) atualizarComunidadeComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	if err := p.garantirDonoTabela(ctx, "comunidades", id, usuarioID, perfilCodigo); err != nil { return comunidadeService.Comunidade{}, err }
	const upd = `UPDATE comunidades SET nome=$2, kind=$3, description=$4, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := p.pool.Exec(ctx, upd, id, corpo.Nome, corpo.Tipo, corpo.Descricao)
	if err != nil { return comunidadeService.Comunidade{}, err }
	if ct.RowsAffected() == 0 { return comunidadeService.Comunidade{}, ErrNaoEncontrado }
	c, ok, err := p.ObterComunidade(ctx, id); if err != nil || !ok { return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade") }
	return c, nil
}
func (p *Postgres) removerComunidadeComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "comunidades", id, usuarioID, perfilCodigo); err != nil { return err }
	ct, err := p.pool.Exec(ctx, `DELETE FROM comunidades WHERE id=$1::uuid`, id)
	if err != nil { return err }
	if ct.RowsAffected() == 0 { return ErrNaoEncontrado }
	return nil
}
func (p *Postgres) AtualizarComunidade(ctx context.Context, id, usuarioID string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return p.atualizarComunidadeComPerfil(ctx, id, usuarioID, "padrao", corpo)
}
func (p *Postgres) AtualizarComunidadeComoAdmin(ctx context.Context, id string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return p.atualizarComunidadeComPerfil(ctx, id, "", "sistema_admin", corpo)
}
func (p *Postgres) RemoverComunidade(ctx context.Context, id, usuarioID string) error {
	return p.removerComunidadeComPerfil(ctx, id, usuarioID, "padrao")
}
func (p *Postgres) RemoverComunidadeComoAdmin(ctx context.Context, id string) error {
	return p.removerComunidadeComPerfil(ctx, id, "", "sistema_admin")
}
