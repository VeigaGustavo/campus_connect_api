package database

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	usuarioService "campus_connect_api/internal/modulos/usuario/structs"
)

func (p *Postgres) Autenticar(contexto context.Context, email, senha string) (*UsuarioInterno, error) {
	const sql = `
SELECT u.id::text, u.nome, u.email, pf.codigo
FROM usuarios u
JOIN perfis_usuario pf ON pf.id = u.perfil_id
WHERE lower(trim(u.email)) = lower(trim($1))
  AND u.ativo
  AND u.senha_hash = crypt($2, u.senha_hash)
`
	var usuario UsuarioInterno
	err := p.pool.QueryRow(contexto, sql, email, senha).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.PerfilCodigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, err
	}
	return &usuario, nil
}

func (p *Postgres) CriarUsuario(contexto context.Context, nome, email, senha, perfilCodigo string) (*UsuarioInterno, error) {
	const sql = `
INSERT INTO usuarios (nome, email, senha_hash, perfil_id)
SELECT $1, lower(trim($2)), crypt($3, gen_salt('bf')), pf.id
FROM perfis_usuario pf
WHERE pf.codigo = $4
RETURNING id::text, nome, email
`
	var usuario UsuarioInterno
	err := p.pool.QueryRow(contexto, sql, nome, email, senha, perfilCodigo).Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	if err != nil {
		return nil, err
	}
	usuario.PerfilCodigo = perfilCodigo
	return &usuario, nil
}

func (p *Postgres) CriarUsuarioComCadastro(contexto context.Context, req usuarioService.RequisicaoCadastroUsuario) (*UsuarioInterno, error) {
	perfilCodigo := "padrao"
	switch req.TipoPerfil {
	case "estudante":
		perfilCodigo = "padrao"
	case "comunidade":
		perfilCodigo = "comunidade"
	case "empresa":
		perfilCodigo = "empresa"
	case "universidade":
		perfilCodigo = "universidade"
	}

	tx, err := p.pool.Begin(contexto)
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
	var usuario UsuarioInterno
	if err := tx.QueryRow(contexto, insUsuario,
		req.NomeCompleto, req.Email, req.Senha, req.Cidade, req.Estado, perfilCodigo,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
		return nil, err
	}
	usuario.PerfilCodigo = perfilCodigo

	detalhes := map[string]any{
		"profile_type": req.TipoPerfil,
		"full_name":    req.NomeCompleto, "age": req.Idade, "cpf": req.CPF, "institution": req.Instituicao,
		"city": req.Cidade, "state": req.Estado, "email": req.Email,
		"community_type": req.TipoComunidade, "community_name": req.NomeComunidade,
		"company_name": req.NomeEmpresa, "company_cnpj": req.CNPJ, "company_description": req.DescricaoEmpresa,
		"institution_name": req.NomeInstituicao, "institution_acronym": req.SiglaInstituicao,
		"institution_type": req.TipoInstituicao, "institution_description": req.DescricaoInstituicao,
	}
	js, _ := json.Marshal(detalhes)
	_, err = tx.Exec(contexto, `
INSERT INTO cadastros_usuario (usuario_id, profile_type, details_json)
VALUES ($1::uuid, $2, $3::jsonb)
`, usuario.ID, req.TipoPerfil, js)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(contexto); err != nil {
		return nil, err
	}
	return &usuario, nil
}
