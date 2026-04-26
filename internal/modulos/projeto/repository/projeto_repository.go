package repository

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	projetoService "campus_connect_api/internal/modulos/projeto/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type projetoRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoProjetoRepository(pool *pgxpool.Pool) projetoService.ProjetoRepository {
	return &projetoRepositoryPostgres{pool: pool}
}

func (repositorio *projetoRepositoryPostgres) ListarProjetos(contexto context.Context) ([]projetoService.Projeto, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []projetoService.Projeto
	for rows.Next() {
		var pr projetoService.Projeto
		if err := rows.Scan(&pr.Identificador, &pr.Titulo, &pr.Descricao); err != nil {
			return nil, err
		}
		out = append(out, pr)
	}
	return out, rows.Err()
}

func (repositorio *projetoRepositoryPostgres) InserirProjeto(contexto context.Context, criadoPor string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const ins = `INSERT INTO projetos (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(contexto, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return projetoService.Projeto{}, err
	}
	if err := inserirCartaoFeedTx(contexto, tx, "project", "dsc-"+id, corpo.Titulo, "Projeto", corpo.Descricao, "Projeto", id, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return projetoService.Projeto{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return projetoService.Projeto{}, err
	}
	pr, ok, err := repositorio.obterProjeto(contexto, id)
	if err != nil || !ok {
		return projetoService.Projeto{}, errors.New("falha ao recarregar projeto")
	}
	return pr, nil
}

func (repositorio *projetoRepositoryPostgres) AtualizarProjeto(contexto context.Context, id, usuarioID string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return repositorio.atualizarProjetoComPerfil(contexto, id, usuarioID, "padrao", corpo)
}

func (repositorio *projetoRepositoryPostgres) AtualizarProjetoComoAdmin(contexto context.Context, id string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	return repositorio.atualizarProjetoComPerfil(contexto, id, "", "sistema_admin", corpo)
}

func (repositorio *projetoRepositoryPostgres) RemoverProjeto(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerProjetoComPerfil(contexto, id, usuarioID, "padrao")
}

func (repositorio *projetoRepositoryPostgres) RemoverProjetoComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerProjetoComPerfil(contexto, id, "", "sistema_admin")
}

func (repositorio *projetoRepositoryPostgres) obterProjeto(contexto context.Context, id string) (projetoService.Projeto, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos WHERE id=$1::uuid`
	var pr projetoService.Projeto
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&pr.Identificador, &pr.Titulo, &pr.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return projetoService.Projeto{}, false, nil
		}
		return projetoService.Projeto{}, false, err
	}
	return pr, true, nil
}

func (repositorio *projetoRepositoryPostgres) atualizarProjetoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	if err := garantirDonoTabela(contexto, repositorio.pool, "projetos", id, usuarioID, perfilCodigo); err != nil {
		return projetoService.Projeto{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	ct, err := tx.Exec(contexto, `UPDATE projetos SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	if ct.RowsAffected() == 0 {
		return projetoService.Projeto{}, comum.ErrNaoEncontrado
	}
	_, _ = tx.Exec(contexto, `UPDATE feed_cartoes SET titulo=$2, excerpt=$3 WHERE kind='project' AND reference_id=$1`, id, corpo.Titulo, corpo.Descricao)
	if err := tx.Commit(contexto); err != nil {
		return projetoService.Projeto{}, err
	}
	pr, ok, err := repositorio.obterProjeto(contexto, id)
	if err != nil || !ok {
		return projetoService.Projeto{}, errors.New("falha ao recarregar projeto")
	}
	return pr, nil
}

func (repositorio *projetoRepositoryPostgres) removerProjetoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := garantirDonoTabela(contexto, repositorio.pool, "projetos", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='project' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM projetos WHERE id=$1::uuid`, id)
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
	case "projetos":
		q = `SELECT criado_por::text FROM projetos WHERE id=$1::uuid`
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
