package repository

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
	comunidadeService "campus_connect_api/internal/modulos/comunidade/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type comunidadeRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoComunidadeRepository(pool *pgxpool.Pool) comunidadeService.ComunidadeRepository {
	return &comunidadeRepositoryPostgres{pool: pool}
}

func (repositorio *comunidadeRepositoryPostgres) ListarComunidades(contexto context.Context) ([]comunidadeService.Comunidade, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []comunidadeService.Comunidade
	for rows.Next() {
		var c comunidadeService.Comunidade
		if err := rows.Scan(&c.Identificador, &c.Nome, &c.Tipo, &c.Descricao); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (repositorio *comunidadeRepositoryPostgres) InserirComunidade(contexto context.Context, criadoPor string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	const ins = `INSERT INTO comunidades (nome, kind, description, criado_por) VALUES ($1,$2,$3,$4::uuid) RETURNING id::text`
	var id string
	if err := repositorio.pool.QueryRow(contexto, ins, corpo.Nome, corpo.Tipo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return comunidadeService.Comunidade{}, err
	}
	c, ok, err := repositorio.obterComunidade(contexto, id)
	if err != nil || !ok {
		return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade")
	}
	return c, nil
}

func (repositorio *comunidadeRepositoryPostgres) AtualizarComunidade(contexto context.Context, id, usuarioID string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.atualizarComunidadeComPerfil(contexto, id, usuarioID, "padrao", corpo)
}

func (repositorio *comunidadeRepositoryPostgres) AtualizarComunidadeComoAdmin(contexto context.Context, id string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	return repositorio.atualizarComunidadeComPerfil(contexto, id, "", "sistema_admin", corpo)
}

func (repositorio *comunidadeRepositoryPostgres) RemoverComunidade(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerComunidadeComPerfil(contexto, id, usuarioID, "padrao")
}

func (repositorio *comunidadeRepositoryPostgres) RemoverComunidadeComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerComunidadeComPerfil(contexto, id, "", "sistema_admin")
}

func (repositorio *comunidadeRepositoryPostgres) obterComunidade(contexto context.Context, id string) (comunidadeService.Comunidade, bool, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades WHERE id=$1::uuid`
	var c comunidadeService.Comunidade
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&c.Identificador, &c.Nome, &c.Tipo, &c.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return comunidadeService.Comunidade{}, false, nil
		}
		return comunidadeService.Comunidade{}, false, err
	}
	return c, true, nil
}

func (repositorio *comunidadeRepositoryPostgres) atualizarComunidadeComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	if err := garantirDonoTabela(contexto, repositorio.pool, "comunidades", id, usuarioID, perfilCodigo); err != nil {
		return comunidadeService.Comunidade{}, err
	}
	const upd = `UPDATE comunidades SET nome=$2, kind=$3, description=$4, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := repositorio.pool.Exec(contexto, upd, id, corpo.Nome, corpo.Tipo, corpo.Descricao)
	if err != nil {
		return comunidadeService.Comunidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return comunidadeService.Comunidade{}, comum.ErrNaoEncontrado
	}
	c, ok, err := repositorio.obterComunidade(contexto, id)
	if err != nil || !ok {
		return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade")
	}
	return c, nil
}

func (repositorio *comunidadeRepositoryPostgres) removerComunidadeComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := garantirDonoTabela(contexto, repositorio.pool, "comunidades", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	ct, err := repositorio.pool.Exec(contexto, `DELETE FROM comunidades WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return nil
}

func garantirDonoTabela(ctx context.Context, pool *pgxpool.Pool, tabela, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return nil
	}
	var q string
	switch tabela {
	case "comunidades":
		q = `SELECT criado_por::text FROM comunidades WHERE id=$1::uuid`
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
