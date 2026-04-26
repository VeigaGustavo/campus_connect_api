package database

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	comunidadeService "campus_connect_api/internal/modulos/comunidade/structs"
	discoverService "campus_connect_api/internal/modulos/discover/structs"
	eventoService "campus_connect_api/internal/modulos/evento/structs"
	grupoService "campus_connect_api/internal/modulos/grupo/structs"
	leituraService "campus_connect_api/internal/modulos/leitura/structs"
	perfilService "campus_connect_api/internal/modulos/perfil/structs"
	projetoService "campus_connect_api/internal/modulos/projeto/structs"
	universidadeService "campus_connect_api/internal/modulos/universidade/structs"
)

func (p *Postgres) ListarEventos(ctx context.Context) ([]eventoService.EventoCampus, error) {
	const sql = `SELECT id::text, titulo, description, start_at, location, organizer FROM eventos ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []eventoService.EventoCampus
	for rows.Next() {
		var e eventoService.EventoCampus
		var t time.Time
		if err := rows.Scan(&e.Identificador, &e.Titulo, &e.Descricao, &t, &e.Local, &e.Organizador); err != nil {
			return nil, err
		}
		e.InicioEm = t.UTC().Format(time.RFC3339)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (p *Postgres) ObterEvento(ctx context.Context, id string) (eventoService.EventoCampus, bool, error) {
	const sql = `SELECT id::text, titulo, description, start_at, location, organizer FROM eventos WHERE id=$1::uuid`
	var e eventoService.EventoCampus
	var t time.Time
	err := p.pool.QueryRow(ctx, sql, id).Scan(&e.Identificador, &e.Titulo, &e.Descricao, &t, &e.Local, &e.Organizador)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return eventoService.EventoCampus{}, false, nil
		}
		return eventoService.EventoCampus{}, false, err
	}
	e.InicioEm = t.UTC().Format(time.RFC3339)
	return e, true, nil
}

func (p *Postgres) InserirEvento(ctx context.Context, criadoPor string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	t, err := time.Parse(time.RFC3339, corpo.InicioEm)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO eventos (titulo, description, start_at, location, organizer, criado_por) VALUES ($1,$2,$3,$4,$5,$6::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.Descricao, t, corpo.Local, corpo.Organizador, criadoPor).Scan(&id); err != nil {
		return eventoService.EventoCampus{}, err
	}
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "event", cartao, corpo.Titulo, corpo.Local, corpo.Descricao, "Início", corpo.InicioEm, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return eventoService.EventoCampus{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return eventoService.EventoCampus{}, err
	}
	e, ok, err := p.ObterEvento(ctx, id)
	if err != nil || !ok {
		return eventoService.EventoCampus{}, errors.New("falha ao recarregar evento")
	}
	return e, nil
}

func (p *Postgres) AtualizarEvento(ctx context.Context, id, usuarioID, perfilCodigo string, corpo eventoService.RequisicaoEvento) (eventoService.EventoCampus, error) {
	if err := p.garantirDonoTabela(ctx, "eventos", id, usuarioID, perfilCodigo); err != nil {
		return eventoService.EventoCampus{}, err
	}
	t, err := time.Parse(time.RFC3339, corpo.InicioEm)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const upd = `UPDATE eventos SET titulo=$2, description=$3, start_at=$4, location=$5, organizer=$6, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := tx.Exec(ctx, upd, id, corpo.Titulo, corpo.Descricao, t, corpo.Local, corpo.Organizador)
	if err != nil {
		return eventoService.EventoCampus{}, err
	}
	if ct.RowsAffected() == 0 {
		return eventoService.EventoCampus{}, ErrNaoEncontrado
	}
	_, _ = tx.Exec(ctx, `UPDATE feed_cartoes SET titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6 WHERE kind='event' AND reference_id=$1`,
		id, corpo.Titulo, corpo.Local, corpo.Descricao, "Início", corpo.InicioEm)
	if err := tx.Commit(ctx); err != nil {
		return eventoService.EventoCampus{}, err
	}
	e, ok, err := p.ObterEvento(ctx, id)
	if err != nil || !ok {
		return eventoService.EventoCampus{}, errors.New("falha ao recarregar evento")
	}
	return e, nil
}

func (p *Postgres) RemoverEvento(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "eventos", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='event' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM eventos WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) garantirDonoTabela(ctx context.Context, tabela, id, usuarioID, perfil string) error {
	if perfil == "sistema_admin" {
		return nil
	}
	var q string
	switch tabela {
	case "eventos":
		q = `SELECT criado_por::text FROM eventos WHERE id=$1::uuid`
	case "grupos_estudo":
		q = `SELECT criado_por::text FROM grupos_estudo WHERE id=$1::uuid`
	case "comunidades":
		q = `SELECT criado_por::text FROM comunidades WHERE id=$1::uuid`
	case "avisos_universidade":
		q = `SELECT criado_por::text FROM avisos_universidade WHERE id=$1::uuid`
	case "leitura_semanal":
		q = `SELECT criado_por::text FROM leitura_semanal WHERE id=$1::uuid`
	case "projetos":
		q = `SELECT criado_por::text FROM projetos WHERE id=$1::uuid`
	default:
		return ErrProibido
	}
	var dono string
	err := p.pool.QueryRow(ctx, q, id).Scan(&dono)
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

func (p *Postgres) ListarGrupos(ctx context.Context) ([]grupoService.GrupoEstudo, error) {
	const sql = `SELECT id::text, titulo, field_of_study, description, level::text, member_count, schedule_label FROM grupos_estudo ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []grupoService.GrupoEstudo
	for rows.Next() {
		var g grupoService.GrupoEstudo
		var lvl string
		if err := rows.Scan(&g.Identificador, &g.Titulo, &g.AreaEstudo, &g.Descricao, &lvl, &g.TotalMembros, &g.RotuloHorario); err != nil {
			return nil, err
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return grupoService.GrupoEstudo{}, false, nil
		}
		return grupoService.GrupoEstudo{}, false, err
	}
	g.Nivel = grupoService.NivelGrupoEstudo(lvl)
	return g, true, nil
}

func (p *Postgres) InserirGrupo(ctx context.Context, criadoPor string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO grupos_estudo (titulo, field_of_study, description, level, member_count, schedule_label, criado_por)
VALUES ($1,$2,$3,$4::varchar,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, 0, corpo.RotuloHorario, criadoPor).Scan(&id); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "study_group", cartao, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	g, ok, err := p.ObterGrupo(ctx, id)
	if err != nil || !ok {
		return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo")
	}
	return g, nil
}

func (p *Postgres) AtualizarGrupo(ctx context.Context, id, usuarioID, perfilCodigo string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	if err := p.garantirDonoTabela(ctx, "grupos_estudo", id, usuarioID, perfilCodigo); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const upd = `UPDATE grupos_estudo SET titulo=$2, field_of_study=$3, description=$4, level=$5::varchar, schedule_label=$6, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := tx.Exec(ctx, upd, id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, corpo.RotuloHorario)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	if ct.RowsAffected() == 0 {
		return grupoService.GrupoEstudo{}, ErrNaoEncontrado
	}
	_, _ = tx.Exec(ctx, `UPDATE feed_cartoes SET titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6 WHERE kind='study_group' AND reference_id=$1`,
		id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel)
	if err := tx.Commit(ctx); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	g, ok, err := p.ObterGrupo(ctx, id)
	if err != nil || !ok {
		return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo")
	}
	return g, nil
}

func (p *Postgres) RemoverGrupo(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "grupos_estudo", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='study_group' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM grupos_estudo WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) ListarComunidades(ctx context.Context) ([]comunidadeService.Comunidade, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
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

func (p *Postgres) ObterComunidade(ctx context.Context, id string) (comunidadeService.Comunidade, bool, error) {
	const sql = `SELECT id::text, nome, kind, description FROM comunidades WHERE id=$1::uuid`
	var c comunidadeService.Comunidade
	err := p.pool.QueryRow(ctx, sql, id).Scan(&c.Identificador, &c.Nome, &c.Tipo, &c.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return comunidadeService.Comunidade{}, false, nil
		}
		return comunidadeService.Comunidade{}, false, err
	}
	return c, true, nil
}

func (p *Postgres) InserirComunidade(ctx context.Context, criadoPor string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	const ins = `INSERT INTO comunidades (nome, kind, description, criado_por) VALUES ($1,$2,$3,$4::uuid) RETURNING id::text`
	var id string
	if err := p.pool.QueryRow(ctx, ins, corpo.Nome, corpo.Tipo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return comunidadeService.Comunidade{}, err
	}
	c, ok, err := p.ObterComunidade(ctx, id)
	if err != nil || !ok {
		return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade")
	}
	return c, nil
}

func (p *Postgres) AtualizarComunidade(ctx context.Context, id, usuarioID, perfilCodigo string, corpo comunidadeService.RequisicaoCriarComunidade) (comunidadeService.Comunidade, error) {
	if err := p.garantirDonoTabela(ctx, "comunidades", id, usuarioID, perfilCodigo); err != nil {
		return comunidadeService.Comunidade{}, err
	}
	const upd = `UPDATE comunidades SET nome=$2, kind=$3, description=$4, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := p.pool.Exec(ctx, upd, id, corpo.Nome, corpo.Tipo, corpo.Descricao)
	if err != nil {
		return comunidadeService.Comunidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return comunidadeService.Comunidade{}, ErrNaoEncontrado
	}
	c, ok, err := p.ObterComunidade(ctx, id)
	if err != nil || !ok {
		return comunidadeService.Comunidade{}, errors.New("falha ao recarregar comunidade")
	}
	return c, nil
}

func (p *Postgres) RemoverComunidade(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "comunidades", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	ct, err := p.pool.Exec(ctx, `DELETE FROM comunidades WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (p *Postgres) ListarAvisosUniversidade(ctx context.Context) ([]universidadeService.AvisoUniversidade, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
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

func (p *Postgres) ObterAvisoUniversidade(ctx context.Context, id string) (universidadeService.AvisoUniversidade, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM avisos_universidade WHERE id=$1::uuid`
	var a universidadeService.AvisoUniversidade
	err := p.pool.QueryRow(ctx, sql, id).Scan(&a.Identificador, &a.Titulo, &a.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return universidadeService.AvisoUniversidade{}, false, nil
		}
		return universidadeService.AvisoUniversidade{}, false, err
	}
	return a, true, nil
}

func (p *Postgres) InserirAvisoUniversidade(ctx context.Context, criadoPor string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO avisos_universidade (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "notice", cartao, corpo.Titulo, "Aviso", corpo.Descricao, "Tipo", "universidade", id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	a, ok, err := p.ObterAvisoUniversidade(ctx, id)
	if err != nil || !ok {
		return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso")
	}
	return a, nil
}

func (p *Postgres) AtualizarAvisoUniversidade(ctx context.Context, id, usuarioID, perfilCodigo string, corpo universidadeService.RequisicaoCriarAvisoUniversidade) (universidadeService.AvisoUniversidade, error) {
	if err := p.garantirDonoTabela(ctx, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	ct, err := p.pool.Exec(ctx, `UPDATE avisos_universidade SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil {
		return universidadeService.AvisoUniversidade{}, err
	}
	if ct.RowsAffected() == 0 {
		return universidadeService.AvisoUniversidade{}, ErrNaoEncontrado
	}
	a, ok, err := p.ObterAvisoUniversidade(ctx, id)
	if err != nil || !ok {
		return universidadeService.AvisoUniversidade{}, errors.New("falha ao recarregar aviso")
	}
	return a, nil
}

func (p *Postgres) RemoverAvisoUniversidade(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "avisos_universidade", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='notice' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM avisos_universidade WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) ListarLeituraSemanal(ctx context.Context) ([]leituraService.ItemLeituraSemanal, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []leituraService.ItemLeituraSemanal
	for rows.Next() {
		var it leituraService.ItemLeituraSemanal
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta); err != nil {
			return nil, err
		}
		it.Tipo = leituraService.TipoItemLeitura(k)
		out = append(out, it)
	}
	return out, rows.Err()
}

func (p *Postgres) ObterLeituraSemanal(ctx context.Context, id string) (leituraService.ItemLeituraSemanal, bool, error) {
	const sql = `SELECT id::text, kind::text, titulo, source, excerpt, image_url, meta_label FROM leitura_semanal WHERE id=$1::uuid`
	var it leituraService.ItemLeituraSemanal
	var k string
	err := p.pool.QueryRow(ctx, sql, id).Scan(&it.Identificador, &k, &it.Titulo, &it.Fonte, &it.Resumo, &it.URLImagem, &it.RotuloMeta)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return leituraService.ItemLeituraSemanal{}, false, nil
		}
		return leituraService.ItemLeituraSemanal{}, false, err
	}
	it.Tipo = leituraService.TipoItemLeitura(k)
	return it, true, nil
}

func (p *Postgres) InserirLeituraSemanal(ctx context.Context, criadoPor string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO leitura_semanal (kind, titulo, source, excerpt, image_url, meta_label, criado_por) VALUES ($1::varchar,$2,$3,$4,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta, criadoPor).Scan(&id); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "reading", cartao, corpo.Titulo, corpo.Fonte, corpo.Resumo, "Leitura", corpo.Tipo, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	it, ok, err := p.ObterLeituraSemanal(ctx, id)
	if err != nil || !ok {
		return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura")
	}
	return it, nil
}

func (p *Postgres) AtualizarLeituraSemanal(ctx context.Context, id, usuarioID, perfilCodigo string, corpo leituraService.RequisicaoLeituraSemanal) (leituraService.ItemLeituraSemanal, error) {
	if err := p.garantirDonoTabela(ctx, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	ct, err := p.pool.Exec(ctx, `UPDATE leitura_semanal SET kind=$2::varchar, titulo=$3, source=$4, excerpt=$5, image_url=$6, meta_label=$7, atualizado_em=now() WHERE id=$1::uuid`,
		id, corpo.Tipo, corpo.Titulo, corpo.Fonte, corpo.Resumo, corpo.URLImagem, corpo.RotuloMeta)
	if err != nil {
		return leituraService.ItemLeituraSemanal{}, err
	}
	if ct.RowsAffected() == 0 {
		return leituraService.ItemLeituraSemanal{}, ErrNaoEncontrado
	}
	it, ok, err := p.ObterLeituraSemanal(ctx, id)
	if err != nil || !ok {
		return leituraService.ItemLeituraSemanal{}, errors.New("falha ao recarregar leitura")
	}
	return it, nil
}

func (p *Postgres) RemoverLeituraSemanal(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "leitura_semanal", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='reading' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM leitura_semanal WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) ListarProjetos(ctx context.Context) ([]projetoService.Projeto, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos ORDER BY criado_em DESC`
	rows, err := p.pool.Query(ctx, sql)
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

func (p *Postgres) ObterProjeto(ctx context.Context, id string) (projetoService.Projeto, bool, error) {
	const sql = `SELECT id::text, titulo, description FROM projetos WHERE id=$1::uuid`
	var pr projetoService.Projeto
	err := p.pool.QueryRow(ctx, sql, id).Scan(&pr.Identificador, &pr.Titulo, &pr.Descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return projetoService.Projeto{}, false, nil
		}
		return projetoService.Projeto{}, false, err
	}
	return pr, true, nil
}

func (p *Postgres) InserirProjeto(ctx context.Context, criadoPor string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	const ins = `INSERT INTO projetos (titulo, description, criado_por) VALUES ($1,$2,$3::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(ctx, ins, corpo.Titulo, corpo.Descricao, criadoPor).Scan(&id); err != nil {
		return projetoService.Projeto{}, err
	}
	cartao := "dsc-" + id
	if err := p.inserirCartaoFeedTx(ctx, tx, "project", cartao, corpo.Titulo, "Projeto", corpo.Descricao, "Projeto", id, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return projetoService.Projeto{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return projetoService.Projeto{}, err
	}
	pr, ok, err := p.ObterProjeto(ctx, id)
	if err != nil || !ok {
		return projetoService.Projeto{}, errors.New("falha ao recarregar projeto")
	}
	return pr, nil
}

func (p *Postgres) AtualizarProjeto(ctx context.Context, id, usuarioID, perfilCodigo string, corpo projetoService.RequisicaoProjeto) (projetoService.Projeto, error) {
	if err := p.garantirDonoTabela(ctx, "projetos", id, usuarioID, perfilCodigo); err != nil {
		return projetoService.Projeto{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	ct, err := tx.Exec(ctx, `UPDATE projetos SET titulo=$2, description=$3, atualizado_em=now() WHERE id=$1::uuid`, id, corpo.Titulo, corpo.Descricao)
	if err != nil {
		return projetoService.Projeto{}, err
	}
	if ct.RowsAffected() == 0 {
		return projetoService.Projeto{}, ErrNaoEncontrado
	}
	_, _ = tx.Exec(ctx, `UPDATE feed_cartoes SET titulo=$2, excerpt=$3 WHERE kind='project' AND reference_id=$1`, id, corpo.Titulo, corpo.Descricao)
	if err := tx.Commit(ctx); err != nil {
		return projetoService.Projeto{}, err
	}
	pr, ok, err := p.ObterProjeto(ctx, id)
	if err != nil || !ok {
		return projetoService.Projeto{}, errors.New("falha ao recarregar projeto")
	}
	return pr, nil
}

func (p *Postgres) RemoverProjeto(ctx context.Context, id, usuarioID, perfilCodigo string) error {
	if err := p.garantirDonoTabela(ctx, "projetos", id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, _ = tx.Exec(ctx, `DELETE FROM feed_cartoes WHERE kind='project' AND reference_id=$1`, id)
	ct, err := tx.Exec(ctx, `DELETE FROM projetos WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return tx.Commit(ctx)
}

func (p *Postgres) FeedDescobrir(ctx context.Context, filtro string, gruposDoUsuario []string) ([]discoverService.ItemDescobrir, error) {
	f := filtro
	if f == "" {
		f = "all"
	}
	var sql string
	var rows pgx.Rows
	var err error
	switch f {
	case "all":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes ORDER BY criado_em DESC`
	case "internships":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='internship' ORDER BY criado_em DESC`
	case "events":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='event' ORDER BY criado_em DESC`
	case "groups":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='study_group' ORDER BY criado_em DESC`
	case "projects":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='project' ORDER BY criado_em DESC`
	case "readings":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='reading' ORDER BY criado_em DESC`
	case "notices":
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes WHERE kind='notice' ORDER BY criado_em DESC`
	default:
		sql = `SELECT id, kind::text, titulo, subtitle, excerpt, meta_primary, meta_secondary, reference_id FROM feed_cartoes ORDER BY criado_em DESC`
	}
	if len(gruposDoUsuario) == 0 {
		sql = strings.Replace(sql, " FROM feed_cartoes", " FROM feed_cartoes WHERE visibility_scope='all'", 1)
		rows, err = p.pool.Query(ctx, sql)
	} else {
		sql = strings.Replace(sql, " FROM feed_cartoes", " FROM feed_cartoes WHERE (visibility_scope='all' OR (visibility_scope='group' AND visibility_group_id = ANY($1::text[])))", 1)
		rows, err = p.pool.Query(ctx, sql, gruposDoUsuario)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []discoverService.ItemDescobrir
	for rows.Next() {
		var it discoverService.ItemDescobrir
		var k string
		if err := rows.Scan(&it.Identificador, &k, &it.Titulo, &it.Subtitulo, &it.Resumo, &it.MetaPrincipal, &it.MetaSecundaria, &it.IDReferencia); err != nil {
			return nil, err
		}
		it.Categoria = discoverService.CategoriaItemDescobrir(k)
		out = append(out, it)
	}
	return out, rows.Err()
}

func (p *Postgres) PerfilUsuario(ctx context.Context, usuarioID string) (perfilService.PerfilUsuario, error) {
	const sql = `
SELECT nome, coalesce(initials,''), coalesce(cover_image_url,''), coalesce(avatar_image_url,''),
       coalesce(performance_certificate_label,''), coalesce(course_and_semester,''), email, coalesce(city_state,''),
       applications_count, groups_count, events_count,
       coalesce(interests,'[]'::jsonb), coalesce(recent_activity,'[]'::jsonb)
FROM usuarios WHERE id=$1::uuid`
	var u perfilService.PerfilUsuario
	var interestsJSON, recentJSON []byte
	err := p.pool.QueryRow(ctx, sql, usuarioID).Scan(
		&u.Nome, &u.Iniciais, &u.URLImagemCapa, &u.URLImagemAvatar, &u.RotuloCertificadoDesempenho,
		&u.CursoESemestre, &u.Email, &u.CidadeEstado,
		&u.TotalCandidaturas, &u.TotalGrupos, &u.TotalEventos,
		&interestsJSON, &recentJSON,
	)
	if err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	_ = json.Unmarshal(interestsJSON, &u.Interesses)
	_ = json.Unmarshal(recentJSON, &u.AtividadesRecentes)
	if u.Interesses == nil {
		u.Interesses = []perfilService.InteressePerfil{}
	}
	if u.AtividadesRecentes == nil {
		u.AtividadesRecentes = []perfilService.LinhaAtividadePerfil{}
	}
	return u, nil
}
