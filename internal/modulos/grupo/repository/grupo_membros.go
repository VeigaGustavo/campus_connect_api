package repository

import (
	"context"
	"errors"
	"time"

	grupoService "campus_connect_api/internal/modulos/grupo/service"
	"github.com/jackc/pgx/v5"
)

func (repositorio *grupoRepositoryPostgres) ObterVisibilidadeGrupo(contexto context.Context, grupoID string) (string, bool, error) {
	const sql = `SELECT coalesce(nullif(trim(visibility),''),'public') FROM grupos_estudo WHERE id=$1::uuid`
	var vis string
	err := repositorio.pool.QueryRow(contexto, sql, grupoID).Scan(&vis)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, err
	}
	return vis, true, nil
}

func (repositorio *grupoRepositoryPostgres) InserirMembroGrupo(contexto context.Context, grupoID, usuarioID, papel string) error {
	const sql = `INSERT INTO grupos_membros (group_id, user_id, role) VALUES ($1::uuid,$2::uuid,$3)
ON CONFLICT (group_id, user_id) DO NOTHING`
	_, err := repositorio.pool.Exec(contexto, sql, grupoID, usuarioID, papel)
	return err
}

func (repositorio *grupoRepositoryPostgres) InserirMembroGrupoTx(contexto context.Context, tx pgx.Tx, grupoID, usuarioID, papel string) error {
	const sql = `INSERT INTO grupos_membros (group_id, user_id, role) VALUES ($1::uuid,$2::uuid,$3)
ON CONFLICT (group_id, user_id) DO NOTHING`
	_, err := tx.Exec(contexto, sql, grupoID, usuarioID, papel)
	return err
}

func (repositorio *grupoRepositoryPostgres) CriarPedidoEntradaGrupo(contexto context.Context, grupoID, usuarioID string) error {
	const sql = `INSERT INTO grupos_pedidos_entrada (group_id, user_id, status) VALUES ($1::uuid,$2::uuid,'pending')
ON CONFLICT (group_id, user_id) DO UPDATE SET status='pending', criado_em=now()`
	_, err := repositorio.pool.Exec(contexto, sql, grupoID, usuarioID)
	return err
}

func (repositorio *grupoRepositoryPostgres) ListarMembrosGrupo(contexto context.Context, grupoID string) ([]grupoService.MembroGrupo, error) {
	const sql = `SELECT u.id::text, u.nome, coalesce(u.avatar_image_url,''), gm.role, gm.joined_at
FROM grupos_membros gm
JOIN usuarios u ON u.id = gm.user_id
WHERE gm.group_id = $1::uuid
ORDER BY CASE gm.role WHEN 'owner' THEN 0 WHEN 'admin' THEN 1 ELSE 2 END, gm.joined_at ASC`
	rows, err := repositorio.pool.Query(contexto, sql, grupoID)
	if err != nil {
		return repositorio.listarMembrosGrupoFallback(contexto, grupoID)
	}
	defer rows.Close()
	var out []grupoService.MembroGrupo
	for rows.Next() {
		var m grupoService.MembroGrupo
		var entrou time.Time
		if err := rows.Scan(&m.UsuarioID, &m.Nome, &m.URLAvatar, &m.Papel, &entrou); err != nil {
			return nil, err
		}
		m.EntrouEm = entrou.UTC().Format(time.RFC3339)
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(out) > 0 {
		return out, nil
	}
	return repositorio.listarMembrosGrupoFallback(contexto, grupoID)
}

func (repositorio *grupoRepositoryPostgres) listarMembrosGrupoFallback(contexto context.Context, grupoID string) ([]grupoService.MembroGrupo, error) {
	const sql = `SELECT criado_por::text FROM grupos_estudo WHERE id=$1::uuid`
	var donoID string
	if err := repositorio.pool.QueryRow(contexto, sql, grupoID).Scan(&donoID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []grupoService.MembroGrupo{}, nil
		}
		return nil, err
	}
	var m grupoService.MembroGrupo
	const sqlUser = `SELECT id::text, nome, coalesce(avatar_image_url,'') FROM usuarios WHERE id=$1::uuid`
	if err := repositorio.pool.QueryRow(contexto, sqlUser, donoID).Scan(&m.UsuarioID, &m.Nome, &m.URLAvatar); err != nil {
		return []grupoService.MembroGrupo{}, nil
	}
	m.Papel = "owner"
	return []grupoService.MembroGrupo{m}, nil
}

func (repositorio *grupoRepositoryPostgres) incrementarMembrosGrupo(contexto context.Context, grupoID string) {
	_, _ = repositorio.pool.Exec(contexto, `UPDATE grupos_estudo SET member_count = member_count + 1 WHERE id=$1::uuid`, grupoID)
}
