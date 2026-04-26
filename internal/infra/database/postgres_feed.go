package database

import (
	"context"

	"github.com/jackc/pgx/v5"

	empresaService "campus_connect_api/internal/modulos/empresa/structs"
)

func (p *Postgres) inserirCartaoFeedTx(ctx context.Context, tx pgx.Tx, kind, cartaoID, titulo, subtitulo, excerpt, metaPri, metaSec, ref, scope, groupID string) error {
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

func (p *Postgres) atualizarCartaoFeedOportunidadeTx(ctx context.Context, tx pgx.Tx, ref string, corpo empresaService.RequisicaoCriarOportunidade) error {
	const sql = `
UPDATE feed_cartoes SET
  titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6
WHERE kind='internship' AND reference_id=$1`
	_, err := tx.Exec(ctx, sql, ref, corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, "Prazo", corpo.PrazoCandidatura)
	return err
}

func nullSeVazio(s string) any {
	if s == "" {
		return nil
	}
	return s
}
