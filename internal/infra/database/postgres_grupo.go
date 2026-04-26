package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	grupoService "campus_connect_api/internal/modulos/grupo/structs"
)

func (p *Postgres) ListarGrupos(ctx context.Context) ([]grupoService.GrupoEstudo, error) {
	const sql = `SELECT id::text, titulo, field_of_study, description, level::text, member_count, schedule_label FROM grupos_estudo ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []grupoService.GrupoEstudo
	for rows.Next() {
		var g grupoService.GrupoEstudo
		var lvl string
		if err := rows.Scan(&g.Identificador, &g.Titulo, &g.AreaEstudo, &g.Descricao, &lvl, &g.TotalMembros, &g.RotuloHorario); err != nil { return nil, err }
		g.Nivel = grupoService.NivelGrupoEstudo(lvl)
		out = append(out, g)
	}
	return out, rows.Err()
}

func (p *Postgres) ObterGrupo(ctx context.Context, id string) (grupoService.GrupoEstudo, bool, error) {
	const sql = `SELECT id::text, titulo, field_of_study, description, level::text, member_count, schedule_label FROM grupos_estudo WHERE id=$1::uuid`
	var g grupoService.GrupoEstudo
	var lvl string
	err := p.pool.QueryRow(ctx, sql, id).Scan(&g.Identificador, &g.Titulo, &g.AreaEstudo, &g.Descricao, &lvl, &g.TotalMembros, &g.RotuloHorario)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return grupoService.GrupoEstudo{}, false, nil }
		return grupoService.GrupoEstudo{}, false, err
	}
	g.Nivel = grupoService.NivelGrupoEstudo(lvl)
	return g, true, nil
}

func (p *Postgres) InserirGrupo(ctx context.Context, criadoPor string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	tx, err := p.pool.Begin(ctx); if err != nil { return grupoService.GrupoEstudo{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO grupos_estudo (titulo, field_of_study, description, level, member_count, schedule_label, criado_por)
VALUES ($1,$2,$3,$4::varchar,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, 0, corpo.RotuloHorario, criadoPor).Scan(&id); err != nil { return grupoService.GrupoEstudo{}, err }
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "study_group", cartao, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil { return grupoService.GrupoEstudo{}, err }
	if err := tx.Commit(ctx); err != nil { return grupoService.GrupoEstudo{}, err }
	g, ok, err := p.ObterGrupo(ctx, id); if err != nil || !ok { return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo") }
	return g, nil
}

func (p *Postgres) atualizarGrupoComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	if err := p.garantirDonoTabela(ctx, "grupos_estudo", id, usuarioID, perfilCodigo); err != nil { return grupoService.GrupoEstudo{}, err }
	tx, err := p.pool.Begin(ctx); if err != nil { return grupoService.GrupoEstudo{}, err }
	defer func() { _ = tx.Rollback(ctx) }()
	const upd = `UPDATE grupos_estudo SET titulo=$2, field_of_study=$3, description=$4, level=$5::varchar, schedule_label=$6, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := tx.Exec(ctx, upd, id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, corpo.RotuloHorario)
	if err != nil { return grupoService.GrupoEstudo{}, err }
	if ct.RowsAffected() == 0 { return grupoService.GrupoEstudo{}, ErrNaoEncontrado }
	_, _ = tx.Exec(ctx, `UPDATE feed_cartoes SET titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6 WHERE kind='study_group' AND reference_id=$1`,
		id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel)
	if err := tx.Commit(ctx); err != nil { return grupoService.GrupoEstudo{}, err }
	g, ok, err := p.ObterGrupo(ctx, id); if err != nil || !ok { return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo") }
	return g, nil
}
func (p *Postgres) removerGrupoComPerfil(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "grupos_estudo", id, usuarioID, perfilCodigo); err != nil { return err }
	tx, err := p.pool.Begin(ctx); if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='study_group' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM grupos_estudo WHERE id=$1::uuid`, id)
	if err != nil { return err }
	if ct.RowsAffected() == 0 { return ErrNaoEncontrado }
	return tx.Commit(ctx)
}
func (p *Postgres) AtualizarGrupo(ctx context.Context, id, usuarioID string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return p.atualizarGrupoComPerfil(ctx, id, usuarioID, "padrao", corpo)
}
func (p *Postgres) AtualizarGrupoComoAdmin(ctx context.Context, id string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return p.atualizarGrupoComPerfil(ctx, id, "", "sistema_admin", corpo)
}
func (p *Postgres) RemoverGrupo(ctx context.Context, id, usuarioID string) error {
	return p.removerGrupoComPerfil(ctx, id, usuarioID, "padrao")
}
func (p *Postgres) RemoverGrupoComoAdmin(ctx context.Context, id string) error {
	return p.removerGrupoComPerfil(ctx, id, "", "sistema_admin")
}
