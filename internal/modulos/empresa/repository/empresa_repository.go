package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	empresaService "campus_connect_api/internal/modulos/empresa/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type empresaRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoEmpresaRepository(pool *pgxpool.Pool) empresaService.EmpresaRepository {
	return &empresaRepositoryPostgres{pool: pool}
}

func (repositorio *empresaRepositoryPostgres) ListarOportunidades(contexto context.Context) ([]empresaService.Oportunidade, error) {
	const sql = `
SELECT id::text, titulo, company_name, short_description, full_description, apply_deadline,
       work_location::text, type_label, coalesce(requirements, '[]'::jsonb)
FROM oportunidades
ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
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

func (repositorio *empresaRepositoryPostgres) ObterOportunidade(contexto context.Context, id string) (empresaService.Oportunidade, bool, error) {
	const sql = `
SELECT id::text, titulo, company_name, short_description, full_description, apply_deadline,
       work_location::text, type_label, coalesce(requirements, '[]'::jsonb)
FROM oportunidades WHERE id = $1::uuid`
	row := repositorio.pool.QueryRow(contexto, sql, id)
	o, err := scanLinhaOportunidade(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return empresaService.Oportunidade{}, false, nil
		}
		return empresaService.Oportunidade{}, false, err
	}
	return o, true, nil
}

func (repositorio *empresaRepositoryPostgres) InserirOportunidade(contexto context.Context, criadoPor string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	deadline, err := time.Parse(time.RFC3339, corpo.PrazoCandidatura)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	reqs, err := json.Marshal(corpo.Requisitos)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()

	const ins = `
INSERT INTO oportunidades (titulo, company_name, short_description, full_description, apply_deadline, work_location, type_label, requirements, criado_por)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9::uuid)
RETURNING id::text`
	var id string
	err = tx.QueryRow(contexto, ins,
		corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, corpo.DescricaoCompleta,
		deadline, corpo.ModalidadeLocal, corpo.RotuloTipo, reqs, criadoPor,
	).Scan(&id)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := inserirCartaoFeedTx(contexto, tx, "internship", "dsc-"+id, corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, "Prazo", corpo.PrazoCandidatura, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return empresaService.Oportunidade{}, err
	}
	op, ok, err := repositorio.ObterOportunidade(contexto, id)
	if err != nil || !ok {
		return empresaService.Oportunidade{}, errors.New("falha ao recarregar oportunidade")
	}
	return op, nil
}

func (repositorio *empresaRepositoryPostgres) AtualizarOportunidade(contexto context.Context, id, usuarioID string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	if err := repositorio.garantirDonoOportunidade(contexto, id, usuarioID); err != nil {
		return empresaService.Oportunidade{}, err
	}
	return repositorio.atualizarOportunidade(contexto, id, corpo)
}

func (repositorio *empresaRepositoryPostgres) AtualizarOportunidadeComoAdmin(contexto context.Context, id string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	return repositorio.atualizarOportunidade(contexto, id, corpo)
}

func (repositorio *empresaRepositoryPostgres) RemoverOportunidade(contexto context.Context, id, usuarioID string) error {
	if err := repositorio.garantirDonoOportunidade(contexto, id, usuarioID); err != nil {
		return err
	}
	return repositorio.RemoverOportunidadeComoAdmin(contexto, id)
}

func (repositorio *empresaRepositoryPostgres) RemoverOportunidadeComoAdmin(contexto context.Context, id string) error {
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='internship' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM oportunidades WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return tx.Commit(contexto)
}

func (repositorio *empresaRepositoryPostgres) atualizarOportunidade(contexto context.Context, id string, corpo empresaService.RequisicaoCriarOportunidade) (empresaService.Oportunidade, error) {
	deadline, err := time.Parse(time.RFC3339, corpo.PrazoCandidatura)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	reqs, err := json.Marshal(corpo.Requisitos)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()

	const upd = `
UPDATE oportunidades SET
  titulo=$2, company_name=$3, short_description=$4, full_description=$5, apply_deadline=$6,
  work_location=$7, type_label=$8, requirements=$9::jsonb, atualizado_em=now()
WHERE id=$1::uuid`
	ct, err := tx.Exec(contexto, upd, id,
		corpo.Titulo, corpo.NomeEmpresa, corpo.DescricaoCurta, corpo.DescricaoCompleta,
		deadline, corpo.ModalidadeLocal, corpo.RotuloTipo, reqs,
	)
	if err != nil {
		return empresaService.Oportunidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return empresaService.Oportunidade{}, comum.ErrNaoEncontrado
	}
	if err := atualizarCartaoFeedOportunidadeTx(contexto, tx, id, corpo); err != nil {
		return empresaService.Oportunidade{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return empresaService.Oportunidade{}, err
	}
	op, ok, err := repositorio.ObterOportunidade(contexto, id)
	if err != nil || !ok {
		return empresaService.Oportunidade{}, errors.New("falha ao recarregar oportunidade")
	}
	return op, nil
}

func (repositorio *empresaRepositoryPostgres) garantirDonoOportunidade(contexto context.Context, id, usuarioID string) error {
	var dono string
	err := repositorio.pool.QueryRow(contexto, `SELECT criado_por::text FROM oportunidades WHERE id=$1::uuid`, id).Scan(&dono)
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

func scanLinhaOportunidade(row interface{ Scan(dest ...any) error }) (empresaService.Oportunidade, error) {
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

func atualizarCartaoFeedOportunidadeTx(ctx context.Context, tx pgx.Tx, ref string, corpo empresaService.RequisicaoCriarOportunidade) error {
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
