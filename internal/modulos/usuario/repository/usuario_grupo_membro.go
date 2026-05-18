package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (repositorio *usuarioRepositoryPostgres) inserirMembroGrupoTx(contexto context.Context, tx pgx.Tx, grupoID, usuarioID, papel string) error {
	const sql = `INSERT INTO grupos_membros (group_id, user_id, role) VALUES ($1::uuid,$2::uuid,$3)
ON CONFLICT (group_id, user_id) DO NOTHING`
	_, err := tx.Exec(contexto, sql, grupoID, usuarioID, papel)
	return err
}
