package repository

import (
	"context"
	"encoding/json"

	comum "campus_connect_api/internal/modulos/comum"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type usuarioRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoUsuarioRepository(pool *pgxpool.Pool) usuarioService.UsuarioRepository {
	return &usuarioRepositoryPostgres{pool: pool}
}

func (repositorio *usuarioRepositoryPostgres) CriarUsuario(contexto context.Context, nome, email, senha, perfilCodigo string) (*usuarioService.UsuarioInterno, error) {
	const sql = `
INSERT INTO usuarios (nome, email, senha_hash, perfil_id)
SELECT $1, lower(trim($2)), crypt($3, gen_salt('bf')), pf.id
FROM perfis_usuario pf
WHERE pf.codigo = $4
RETURNING id::text, nome, email
`
	var usuario usuarioService.UsuarioInterno
	err := repositorio.pool.QueryRow(contexto, sql, nome, email, senha, perfilCodigo).Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	if err != nil {
		return nil, err
	}
	usuario.PerfilCodigo = perfilCodigo
	return &usuario, nil
}

func (repositorio *usuarioRepositoryPostgres) CriarUsuarioComCadastro(contexto context.Context, requisicao usuarioService.RequisicaoCadastroUsuario) (*usuarioService.UsuarioInterno, error) {
	perfilCodigo := comum.PerfilPadrao
	switch requisicao.TipoPerfil {
	case "estudante":
		perfilCodigo = comum.PerfilPadrao
	case "comunidade":
		perfilCodigo = "comunidade"
	case "empresa":
		perfilCodigo = "empresa"
	case "universidade":
		perfilCodigo = "universidade"
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const insUsuario = `
INSERT INTO usuarios (nome, email, senha_hash, perfil_id, city_state)
SELECT $1, lower(trim($2)), crypt($3, gen_salt('bf')), pf.id, ($4 || ' - ' || $5)
FROM perfis_usuario pf
WHERE pf.codigo = $6
RETURNING id::text, nome, email
`
	var usuario usuarioService.UsuarioInterno
	if err := tx.QueryRow(contexto, insUsuario,
		requisicao.NomeCompleto, requisicao.Email, requisicao.Senha, requisicao.Cidade, requisicao.Estado, perfilCodigo,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
		return nil, err
	}
	usuario.PerfilCodigo = perfilCodigo

	detalhes := map[string]any{
		"profile_type": requisicao.TipoPerfil,
		"full_name":    requisicao.NomeCompleto, "age": requisicao.Idade, "cpf": requisicao.CPF, "institution": requisicao.Instituicao,
		"city": requisicao.Cidade, "state": requisicao.Estado, "email": requisicao.Email,
		"community_type": requisicao.TipoComunidade, "community_name": requisicao.NomeComunidade,
		"company_name": requisicao.NomeEmpresa, "company_cnpj": requisicao.CNPJ, "company_description": requisicao.DescricaoEmpresa,
		"institution_name": requisicao.NomeInstituicao, "institution_acronym": requisicao.SiglaInstituicao,
		"institution_type": requisicao.TipoInstituicao, "institution_description": requisicao.DescricaoInstituicao,
	}
	js, _ := json.Marshal(detalhes)
	_, err = tx.Exec(contexto, `
INSERT INTO cadastros_usuario (usuario_id, profile_type, details_json)
VALUES ($1::uuid, $2, $3::jsonb)
`, usuario.ID, requisicao.TipoPerfil, js)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(contexto); err != nil {
		return nil, err
	}
	return &usuario, nil
}
