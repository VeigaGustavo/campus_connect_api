package repository

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	universidadeService "campus_connect_api/internal/modulos/universidade/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type universidadeRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoUniversidadeRepository(pool *pgxpool.Pool) universidadeService.UniversidadeRepository {
	return &universidadeRepositoryPostgres{pool: pool}
}

func (repositorio *universidadeRepositoryPostgres) ListarAvisosUniversidade(contexto context.Context) ([]universidadeService.AvisoUniversidade, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []universidadeService.AvisoUniversidade
	for rows.Next() {
		var a universidadeService.AvisoUniversidade
		if err := rows.Scan(&a.Identificador, &a.Titulo, &a.Descricao); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (repositorio *universidadeRepositoryPostgres) InserirAvisoUniversidade(contexto context.Context, criadoPor string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const ins = `INSERT INTO avisos_universidade (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(contexto, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	if err := inserirCartaoFeedTx(contexto, tx, "notice", "dsc-"+id, corpo.Titulo, "Aviso", corpo.Descricao, "Tipo", "universidade", id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	a, ok, err := repositorio.obterAvisoUniversidade(contexto, id)
	if err != nil || !ok {
		return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso")
	}
	return a, nil
}

func (repositorio *universidadeRepositoryPostgres) AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return repositorio.atualizarAvisoUniversidadeComPerfil(contexto, id, usuarioID, "padrao", corpo)
}

func (repositorio *universidadeRepositoryPostgres) AtualizarAvisoUniversidadeComoAdmin(contexto context.Context, id string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	return repositorio.atualizarAvisoUniversidadeComPerfil(contexto, id, "", "sistema_admin", corpo)
}

func (repositorio *universidadeRepositoryPostgres) RemoverAvisoUniversidade(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerAvisoUniversidadeComPerfil(contexto, id, usuarioID, "padrao")
}

func (repositorio *universidadeRepositoryPostgres) RemoverAvisoUniversidadeComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerAvisoUniversidadeComPerfil(contexto, id, "", "sistema_admin")
}

func (repositorio *universidadeRepositoryPostgres) obterAvisoUniversidade(contexto context.Context, id string) (universidadeService.AvisoUniversidade, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade WHERE id=$1::uuid`
	var a universidadeService.AvisoUniversidade
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&a.Identificador, &a.Titulo, &a.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return universidadeService.AvisoUniversidade{}, false, nil
		}
		return universidadeService.AvisoUniversidade{}, false, err
	}
	return a, true, nil
}

func (repositorio *universidadeRepositoryPostgres) atualizarAvisoUniversidadeComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	if err := garantirDonoTabela(contexto, repositorio.pool, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	ct, err := repositorio.pool.Exec(contexto, `UPDATE avisos_universidade SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return universidadeService.AvisoUniversidade{}, comum.ErrNaoEncontrado
	}
	a, ok, err := repositorio.obterAvisoUniversidade(contexto, id)
	if err != nil || !ok {
		return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso")
	}
	return a, nil
}

func (repositorio *universidadeRepositoryPostgres) removerAvisoUniversidadeComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := garantirDonoTabela(contexto, repositorio.pool, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='notice' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM avisos_universidade WHERE id=$1::uuid`, id)
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
	case "avisos_universidade":
		q = `SELECT criado_por::text FROM avisos_universidade WHERE id=$1::uuid`
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
