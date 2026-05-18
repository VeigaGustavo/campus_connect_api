package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("perfil %q inexistente em perfis_usuario (INSERT nao encontrou linha); inclua esse codigo na tabela antes do cadastro", perfilCodigo)
		}
		return nil, err
	}
	usuario.PerfilCodigo = perfilCodigo
	return &usuario, nil
}

func (repositorio *usuarioRepositoryPostgres) CriarUsuarioComCadastro(contexto context.Context, requisicao usuarioService.RequisicaoCadastroUsuario) (*usuarioService.UsuarioInterno, *usuarioService.ResultadoCadastroComunidade, error) {
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
		return nil, nil, err
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("perfil %q inexistente em perfis_usuario (INSERT nao encontrou linha); inclua esse codigo na tabela antes do cadastro", perfilCodigo)
		}
		return nil, nil, err
	}
	usuario.PerfilCodigo = perfilCodigo

	detalhes := map[string]any{
		"profile_type": requisicao.TipoPerfil,
		"full_name":    requisicao.NomeCompleto, "age": requisicao.Idade, "birth_date": strings.TrimSpace(requisicao.DataNascimento),
		"cpf": requisicao.CPF, "institution": requisicao.Instituicao,
		"city": requisicao.Cidade, "state": requisicao.Estado, "email": requisicao.Email,
		"community_type": requisicao.TipoComunidade, "community_name": requisicao.NomeComunidade,
		"group_title": requisicao.TituloGrupo, "group_description": requisicao.DescricaoGrupo,
		"group_visibility": requisicao.VisibilidadeGrupo,
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
		return nil, nil, err
	}

	var resultadoComunidade *usuarioService.ResultadoCadastroComunidade
	if requisicao.TipoPerfil == "comunidade" {
		res, err := repositorio.criarComunidadeEGrupoNoCadastro(contexto, tx, usuario.ID, requisicao)
		if err != nil {
			return nil, nil, err
		}
		resultadoComunidade = &res
	}

	if err := tx.Commit(contexto); err != nil {
		return nil, nil, err
	}
	return &usuario, resultadoComunidade, nil
}

func (repositorio *usuarioRepositoryPostgres) criarComunidadeEGrupoNoCadastro(
	contexto context.Context,
	tx pgx.Tx,
	usuarioID string,
	requisicao usuarioService.RequisicaoCadastroUsuario,
) (usuarioService.ResultadoCadastroComunidade, error) {
	tituloGrupo := strings.TrimSpace(requisicao.TituloGrupo)
	if tituloGrupo == "" {
		tituloGrupo = strings.TrimSpace(requisicao.NomeComunidade)
	}
	descricao := strings.TrimSpace(requisicao.DescricaoGrupo)
	visibilidade := strings.TrimSpace(requisicao.VisibilidadeGrupo)
	if visibilidade != "private" {
		visibilidade = "public"
	}
	tipoComunidade := strings.TrimSpace(requisicao.TipoComunidade)
	rotuloArea := tipoComunidade
	switch tipoComunidade {
	case "atletica":
		rotuloArea = "Atlética"
	case "ca":
		rotuloArea = "Centro Acadêmico"
	}

	const insComunidade = `
INSERT INTO comunidades (nome, kind, description, criado_por)
VALUES ($1, $2, $3, $4::uuid)
RETURNING id::text`
	var comunidadeID string
	if err := tx.QueryRow(contexto, insComunidade,
		strings.TrimSpace(requisicao.NomeComunidade),
		tipoComunidade,
		descricao,
		usuarioID,
	).Scan(&comunidadeID); err != nil {
		return usuarioService.ResultadoCadastroComunidade{}, err
	}

	const insGrupo = `
INSERT INTO grupos_estudo (titulo, field_of_study, description, level, member_count, schedule_label, criado_por, visibility)
VALUES ($1, $2, $3, 'beginner', 1, '', $4::uuid, $5)
RETURNING id::text`
	var grupoID string
	if err := tx.QueryRow(contexto, insGrupo, tituloGrupo, rotuloArea, descricao, usuarioID, visibilidade).Scan(&grupoID); err != nil {
		const insGrupoLegado = `
INSERT INTO grupos_estudo (titulo, field_of_study, description, level, member_count, schedule_label, criado_por)
VALUES ($1, $2, $3, 'beginner', 1, $5, $4::uuid)
RETURNING id::text`
		if errLegado := tx.QueryRow(contexto, insGrupoLegado, tituloGrupo, rotuloArea, descricao, usuarioID, "").Scan(&grupoID); errLegado != nil {
			return usuarioService.ResultadoCadastroComunidade{}, errLegado
		}
	}

	if err := repositoryutil.InserirCartaoFeedTx(
		contexto, tx, comum.FeedKindGrupoEstudo, "dsc-"+grupoID,
		tituloGrupo, rotuloArea, descricao, "Nível", "beginner", grupoID, "", "",
	); err != nil {
		return usuarioService.ResultadoCadastroComunidade{}, err
	}
	_ = repositorio.inserirMembroGrupoTx(contexto, tx, grupoID, usuarioID, "owner")
	return usuarioService.ResultadoCadastroComunidade{
		CommunityID: comunidadeID,
		GroupID:     grupoID,
	}, nil
}
