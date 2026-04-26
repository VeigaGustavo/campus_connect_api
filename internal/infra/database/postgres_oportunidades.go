package database

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	empresaService "campus_connect_api/internal/modulos/empresa/structs"
)

type scannerLinha interface {
	Scan(dest ...any) error
}

func (p *Postgres) ListarOportunidades(ctx context.Context) ([]empresaService.Oportunidade, error) {
	const sql = `
SELECT id::text, titulo, company_name, short_description, full_description, apply_deadline,
       work_location::text, type_label, coalesce(requirements, '[]'::jsonb)
FROM oportunidades
ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []empresaService.Oportunidade
	for rows.Next() {
		o, err := scanLinhaOportunidade(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (p *Postgres) ObterOportunidade(ctx context.Context, id string) (empresaService.Oportunidade, bool, error) {
	const sql = `
SELECT id::text, titulo, company_name, short_description, full_description, apply_deadline,
       work_location::text, type_label, coalesce(requirements, '[]'::jsonb)
FROM oportunidades WHERE id = $1::uuid`
	row := p.pool.QueryRow(ctx, sql, id)
	o, err := scanLinhaOportunidade(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return empresaService.Oportunidade{}, false, nil
		}
		return empresaService.Oportunidade{}, false, err
	}
	return o, true, nil
}

func scanLinhaOportunidade(row scannerLinha) (empresaService.Oportunidade, error) {
	var o empresaService.Oportunidade
	var deadline time.Time
	var wl string
	var reqsJSON []byte
	err := row.Scan(
		&o.Identificador, &o.Titulo, &o.NomeEmpresa, &o.DescricaoCurta, &o.DescricaoCompleta,
		&deadline, &wl, &o.RotuloTipo, &reqsJSON,
	)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	o.PrazoCandidatura = deadline.UTC().Format(time.RFC3339)
	o.ModalidadeLocal = empresaService.ModalidadeLocalTrabalho(wl)
	if len(reqsJSON) > 0 {
		_ = json.Unmarshal(reqsJSON, &o.Requisitos)
	}
	if o.Requisitos == nil {
		o.Requisitos = []string{}
	}
	return o, nil
}

func (p *Postgres) InserirOportunidade(ctx context.Context, criadoPor string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	deadline, err := time.Parse(time.RFC3339, corpo.PrazoCandidatura)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	reqs, err := json.Marshal(corpo.Requisitos)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const ins = `
INSERT INTO oportunidades (titulo, company_name, short_description, full_description, apply_deadline, work_location, type_label, requirements, criado_por)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9::uuid)
RETURNING id::text`
	var id string
	err = tx.QueryRow(ctx, ins,
		corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, corpo.DescricaoCompleta,
		deadline, corpo.ModalidadeLocal, corpo.RotuloTipo, reqs, criadoPor,
	).Scan(&id)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	cartaoID := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "internship", cartaoID, corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, "Prazo", corpo.PrazoCandidatura, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return empresaService.Oportunidade{}, err
	}
	op, ok, err := p.ObterOportunidade(ctx, id)
	if err != nil || !ok {
		return empresaService.Oportunidade{}, errors.New("falha ao recarregar oportunidade")
	}
	return op, nil
}

func (p *Postgres) AtualizarOportunidade(ctx context.Context, id, usuarioID string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	if err := p.garantirDonoOportunidade(ctx, id, usuarioID); err != nil {
		return empresaService.Oportunidade{}, err
	}
	deadline, err := time.Parse(time.RFC3339, corpo.PrazoCandidatura)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	reqs, err := json.Marshal(corpo.Requisitos)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const upd = `
UPDATE oportunidades SET
  titulo=$2, company_name=$3, short_description=$4, full_description=$5, apply_deadline=$6,
  work_location=$7, type_label=$8, requirements=$9::jsonb, atualizado_em=now()
WHERE id=$1::uuid`
	ct, err := tx.Exec(ctx, upd, id,
		corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, corpo.DescricaoCompleta,
		deadline, corpo.ModalidadeLocal, corpo.RotuloTipo, reqs,
	)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return empresaService.Oportunidade{}, ErrNaoEncontrado
	}
	if err := p.atualizarCartaoFeedOportunidadeTx(ctx, tx, id, corpo); err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return empresaService.Oportunidade{}, err
	}
	op, ok, err := p.ObterOportunidade(ctx, id)
	if err != nil || !ok {
		return empresaService.Oportunidade{}, errors.New("falha ao recarregar oportunidade")
	}
	return op, nil
}

func (p *Postgres) AtualizarOportunidadeComoAdmin(ctx context.Context, id string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	deadline, err := time.Parse(time.RFC3339, corpo.PrazoCandidatura)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	reqs, err := json.Marshal(corpo.Requisitos)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const upd = `
UPDATE oportunidades SET
  titulo=$2, company_name=$3, short_description=$4, full_description=$5, apply_deadline=$6,
  work_location=$7, type_label=$8, requirements=$9::jsonb, atualizado_em=now()
WHERE id=$1::uuid`
	ct, err := tx.Exec(ctx, upd, id,
		corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, corpo.DescricaoCompleta,
		deadline, corpo.ModalidadeLocal, corpo.RotuloTipo, reqs,
	)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return empresaService.Oportunidade{}, ErrNaoEncontrado
	}
	if err := p.atualizarCartaoFeedOportunidadeTx(ctx, tx, id, corpo); err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return empresaService.Oportunidade{}, err
	}
	op, ok, err := p.ObterOportunidade(ctx, id)
	if err != nil || !ok {
		return empresaService.Oportunidade{}, errors.New("falha ao recarregar oportunidade")
	}
	return op, nil
}

func (p *Postgres) RemoverOportunidade(ctx context.Context, id, usuarioID string) error {
	if err := p.garantirDonoOportunidade(ctx, id, usuarioID); err != nil {
		return err
	}
	return p.RemoverOportunidadeComoAdmin(ctx, id)
}

func (p *Postgres) RemoverOportunidadeComoAdmin(ctx context.Context, id string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='internship' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM oportunidades WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) garantirDonoOportunidade(ctx context.Context, id, usuarioID string) error {
	const sql = `SELECT criado_por::text FROM oportunidades WHERE id=$1::uuid`
	var dono string
	err := p.pool.QueryRow(ctx, sql, id).Scan(&dono)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNaoEncontrado
		}
		return err
	}
	if dono != usuarioID {
		return ErrProibido
	}
	return nil
}
