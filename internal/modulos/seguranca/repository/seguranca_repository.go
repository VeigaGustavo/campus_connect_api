package repository

import (
	"context"
	"errors"

	segurancaService "campus_connect_api/internal/modulos/seguranca/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type segurancaRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoSegurancaRepository(pool *pgxpool.Pool) segurancaService.SegurancaRepository {
	return &segurancaRepositoryPostgres{pool: pool}
}

func (repositorio *segurancaRepositoryPostgres) Autenticar(contexto context.Context, email, senha string) (*segurancaService.UsuarioAutenticado, error) {
	const sql = `
SELECT u.id::text, u.nome, u.email, pf.codigo
FROM usuarios u
JOIN perfis_usuario pf ON pf.id = u.perfil_id
WHERE lower(trim(u.email)) = lower(trim($1))
  AND u.ativo
  AND u.senha_hash = crypt($2, u.senha_hash)
`
	var usuario segurancaService.UsuarioAutenticado
	err := repositorio.pool.QueryRow(contexto, sql, email, senha).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.PerfilCodigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, segurancaService.ErrAutenticacaoInvalida
		}
		return nil, err
	}
	return &usuario, nil
}
