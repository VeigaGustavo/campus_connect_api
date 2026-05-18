package repository

import (
	"context"
	"errors"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
	eventoService "campus_connect_api/internal/modulos/evento/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type eventoRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoEventoRepository(pool *pgxpool.Pool) eventoService.EventoRepository {
	return &eventoRepositoryPostgres{pool: pool}
}

func (repositorio *eventoRepositoryPostgres) ListarEventos(contexto context.Context) ([]eventoService.EventoCampus, error) {
	const sql = `SELECT id::text, titulo, description, start_at, location, organizer, criado_por::text FROM eventos ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []eventoService.EventoCampus
	for rows.Next() {
		var e eventoService.EventoCampus
		var t time.Time
		if err := rows.Scan(&e.Identificador, &e.Titulo, &e.Descricao, &t, &e.Local, &e.Organizador, &e.AutorID); err != nil {
			return nil, err
		}
		if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, e.AutorID, &e.Autor); err != nil {
			return nil, err
		}
		e.InicioEm = t.UTC().Format(time.RFC3339)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (repositorio *eventoRepositoryPostgres) InserirEvento(contexto context.Context, criadoPor string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	t, err := time.Parse(time.RFC3339, corpo.InicioEm)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const ins = `INSERT INTO eventos (titulo, description, start_at, location, organizer, criado_por) VALUES ($1,$2,$3,$4,$5,$6::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(contexto, ins, corpo.Titulo, corpo.Descricao, t, corpo.Local, corpo.Organizador, criadoPor).Scan(&id); err != nil {
		return eventoService.EventoCampus{}, err
	}
	if err := repositoryutil.InserirCartaoFeedTx(contexto, tx, comum.FeedKindEvento, "dsc-"+id, corpo.Titulo, corpo.Local, corpo.Descricao, "Início", corpo.InicioEm, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return eventoService.EventoCampus{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return eventoService.EventoCampus{}, err
	}
	e, ok, err := repositorio.obterEvento(contexto, id)
	if err != nil || !ok {
		return eventoService.EventoCampus{}, errors.New("falha ao recarregar evento")
	}
	return e, nil
}

func (repositorio *eventoRepositoryPostgres) AtualizarEvento(contexto context.Context, id, usuarioID string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.atualizarEventoComPerfil(contexto, id, usuarioID, comum.PerfilPadrao, corpo)
}

func (repositorio *eventoRepositoryPostgres) AtualizarEventoComoAdmin(contexto context.Context, id string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	return repositorio.atualizarEventoComPerfil(contexto, id, "", comum.PerfilSistemaAdmin, corpo)
}

func (repositorio *eventoRepositoryPostgres) RemoverEvento(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerEventoComPerfil(contexto, id, usuarioID, comum.PerfilPadrao)
}

func (repositorio *eventoRepositoryPostgres) RemoverEventoComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerEventoComPerfil(contexto, id, "", comum.PerfilSistemaAdmin)
}

func (repositorio *eventoRepositoryPostgres) obterEvento(contexto context.Context, id string) (eventoService.EventoCampus, bool, error) {
	const sql = `SELECT id::text, titulo, description, start_at, location, organizer, criado_por::text FROM eventos WHERE id=$1::uuid`
	var e eventoService.EventoCampus
	var t time.Time
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&e.Identificador, &e.Titulo, &e.Descricao, &t, &e.Local, &e.Organizador, &e.AutorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return eventoService.EventoCampus{}, false, nil
		}
		return eventoService.EventoCampus{}, false, err
	}
	if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, e.AutorID, &e.Autor); err != nil {
		return eventoService.EventoCampus{}, false, err
	}
	e.InicioEm = t.UTC().Format(time.RFC3339)
	return e, true, nil
}

func (repositorio *eventoRepositoryPostgres) atualizarEventoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM eventos WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
		return eventoService.EventoCampus{}, err
	}
	t, err := time.Parse(time.RFC3339, corpo.InicioEm)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const upd = `UPDATE eventos SET titulo=$2, description=$3, start_at=$4, location=$5, organizer=$6, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := tx.Exec(contexto, upd, id, corpo.Titulo, corpo.Descricao, t, corpo.Local, corpo.Organizador)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	if ct.RowsAffected() == 0 {
		return eventoService.EventoCampus{}, comum.ErrNaoEncontrado
	}
	_, _ = tx.Exec(contexto, `UPDATE feed_cartoes SET titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6 WHERE kind='event' AND reference_id=$1`,
		id, corpo.Titulo, corpo.Local, corpo.Descricao, "Início", corpo.InicioEm)
	if err := tx.Commit(contexto); err != nil {
		return eventoService.EventoCampus{}, err
	}
	e, ok, err := repositorio.obterEvento(contexto, id)
	if err != nil || !ok {
		return eventoService.EventoCampus{}, errors.New("falha ao recarregar evento")
	}
	return e, nil
}

func (repositorio *eventoRepositoryPostgres) removerEventoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM eventos WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='event' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM eventos WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return tx.Commit(contexto)
}

