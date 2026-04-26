package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

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
