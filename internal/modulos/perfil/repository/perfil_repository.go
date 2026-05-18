package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type perfilRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NovoPerfilRepository(pool *pgxpool.Pool) perfilService.PerfilRepository {
	return &perfilRepositoryPostgres{pool: pool}
}

func (repositorio *perfilRepositoryPostgres) PerfilUsuario(contexto context.Context, usuarioID string) (perfilService.PerfilUsuario, error) {
	const sql = `
SELECT id::text, nome, coalesce(cover_image_url,''), coalesce(avatar_image_url,''),
       email, coalesce(city_state,''),
       coalesce(about_me,''), coalesce(job_title,''), coalesce(course,''), coalesce(semester,''), coalesce(institution_name,''),
       applications_count, groups_count, events_count,
       coalesce(interests,'[]'::jsonb), coalesce(favorite_topics,'[]'::jsonb), coalesce(specialties,'[]'::jsonb)
FROM usuarios WHERE id=$1::uuid`
	var u perfilService.PerfilUsuario
	var interessesJSON, topicosJSON, especialidadesJSON []byte
	err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(
		&u.Identificador, &u.Nome, &u.URLImagemCapa, &u.URLImagemAvatar,
		&u.Email, &u.CidadeEstado,
		&u.SobreMim, &u.Cargo, &u.Curso, &u.Semestre, &u.Instituicao,
		&u.TotalCandidaturas, &u.TotalGrupos, &u.TotalEventos,
		&interessesJSON, &topicosJSON, &especialidadesJSON,
	)
	if err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	_ = json.Unmarshal(interessesJSON, &u.Interesses)
	_ = json.Unmarshal(topicosJSON, &u.TopicosFavoritos)
	_ = json.Unmarshal(especialidadesJSON, &u.Especialidades)
	if u.Interesses == nil {
		u.Interesses = []string{}
	}
	if u.TopicosFavoritos == nil {
		u.TopicosFavoritos = []string{}
	}
	if u.Especialidades == nil {
		u.Especialidades = []string{}
	}
	destaque, err := repositorio.obterDestaqueComunidade(contexto, usuarioID)
	if err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	u.DestaqueComunidade = destaque
	u.ContextoPerfil = "user"
	repositorio.aplicarContratoFrontPerfil(contexto, &u, usuarioID, repositorio.obterPapelContaUsuario(contexto, usuarioID))
	return u, nil
}

func (repositorio *perfilRepositoryPostgres) PerfilParaExibicao(contexto context.Context, usuarioID, perfilCodigoConta string) (perfilService.PerfilUsuario, error) {
	switch perfilCodigoConta {
	case "comunidade", "empresa", "universidade":
		u, err := repositorio.perfilInstitucional(contexto, usuarioID, perfilCodigoConta)
		if err != nil {
			return perfilService.PerfilUsuario{}, err
		}
		if u.ContextoPerfil == "" {
			u.ContextoPerfil = "organization"
		}
		return u, nil
	default:
		return repositorio.PerfilUsuario(contexto, usuarioID)
	}
}

const limItensPainelOrganizacao = 12

func (repositorio *perfilRepositoryPostgres) enriquecerPainelOrganizacao(contexto context.Context, u *perfilService.PerfilUsuario, usuarioID, perfilCodigoConta string, detalhes map[string]any) error {
	p := &perfilService.PainelOrganizacaoPerfil{}
	posts, nPosts, err := repositorio.listarPostsOrganizacao(contexto, usuarioID, limItensPainelOrganizacao)
	if err != nil {
		return err
	}
	p.Posts, p.PostsTotal = posts, nPosts

	evs, nEv, err := repositorio.listarEventosOrganizacao(contexto, usuarioID, limItensPainelOrganizacao)
	if err != nil {
		return err
	}
	p.Eventos, p.EventosTotal = evs, nEv

	switch perfilCodigoConta {
	case "empresa":
		jobs, nJ, err := repositorio.listarOportunidadesOrganizacao(contexto, usuarioID, limItensPainelOrganizacao)
		if err != nil {
			return err
		}
		p.Vagas, p.VagasTotal = jobs, nJ
	case "universidade":
		p.MapURL = textoEmMapa(detalhes, "map_url")
	case "comunidade":
		p.InstituicaoPai = textoEmMapa(detalhes, "institution")
		gr, nG, err := repositorio.listarGruposOrganizacao(contexto, usuarioID, limItensPainelOrganizacao)
		if err != nil {
			return err
		}
		p.Grupos, p.GruposTotal = gr, nG
	}
	u.PainelOrganizacao = p
	return nil
}

func (repositorio *perfilRepositoryPostgres) listarPostsOrganizacao(contexto context.Context, usuarioID string, lim int) ([]perfilService.ResumoPostOrganizacao, int, error) {
	var total int
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM feed_posts WHERE author_id=$1::uuid`, usuarioID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repositorio.pool.Query(contexto, `
SELECT id::text, coalesce(left(body_text, 280), ''), criado_em
FROM feed_posts WHERE author_id=$1::uuid ORDER BY criado_em DESC LIMIT $2`, usuarioID, lim)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []perfilService.ResumoPostOrganizacao
	for rows.Next() {
		var it perfilService.ResumoPostOrganizacao
		var criadoEm time.Time
		if err := rows.Scan(&it.Identificador, &it.Antevisao, &criadoEm); err != nil {
			return nil, 0, err
		}
		it.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
		out = append(out, it)
	}
	return out, total, rows.Err()
}

func (repositorio *perfilRepositoryPostgres) listarEventosOrganizacao(contexto context.Context, usuarioID string, lim int) ([]perfilService.ResumoPublicacaoOrganizacao, int, error) {
	var total int
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM eventos WHERE criado_por=$1::uuid`, usuarioID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repositorio.pool.Query(contexto, `
SELECT id::text, titulo, coalesce(nullif(trim(location),''), nullif(trim(organizer),''), ''), criado_em
FROM eventos WHERE criado_por=$1::uuid ORDER BY criado_em DESC LIMIT $2`, usuarioID, lim)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return escanearResumosPublicacao(rows, total)
}

func (repositorio *perfilRepositoryPostgres) listarOportunidadesOrganizacao(contexto context.Context, usuarioID string, lim int) ([]perfilService.ResumoPublicacaoOrganizacao, int, error) {
	var total int
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM oportunidades WHERE criado_por=$1::uuid`, usuarioID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repositorio.pool.Query(contexto, `
SELECT id::text, titulo, coalesce(company_name,''), criado_em
FROM oportunidades WHERE criado_por=$1::uuid ORDER BY criado_em DESC LIMIT $2`, usuarioID, lim)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return escanearResumosPublicacao(rows, total)
}

func (repositorio *perfilRepositoryPostgres) listarGruposOrganizacao(contexto context.Context, usuarioID string, lim int) ([]perfilService.ResumoPublicacaoOrganizacao, int, error) {
	var total int
	if err := repositorio.pool.QueryRow(contexto, `SELECT count(1) FROM grupos_estudo WHERE criado_por=$1::uuid`, usuarioID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repositorio.pool.Query(contexto, `
SELECT id::text, titulo, coalesce(field_of_study,''), criado_em
FROM grupos_estudo WHERE criado_por=$1::uuid ORDER BY criado_em DESC LIMIT $2`, usuarioID, lim)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return escanearResumosPublicacao(rows, total)
}

func escanearResumosPublicacao(rows pgx.Rows, total int) ([]perfilService.ResumoPublicacaoOrganizacao, int, error) {
	var out []perfilService.ResumoPublicacaoOrganizacao
	for rows.Next() {
		var it perfilService.ResumoPublicacaoOrganizacao
		var criadoEm time.Time
		if err := rows.Scan(&it.Identificador, &it.Titulo, &it.Subtitulo, &criadoEm); err != nil {
			return nil, 0, err
		}
		it.IDReferencia = it.Identificador
		it.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
		out = append(out, it)
	}
	return out, total, rows.Err()
}

func (repositorio *perfilRepositoryPostgres) perfilInstitucional(contexto context.Context, usuarioID, perfilCodigoConta string) (perfilService.PerfilUsuario, error) {
	detalhes, err := repositorio.carregarDetalhesCadastro(contexto, usuarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repositorio.perfilInstitucionalSemCadastro(contexto, usuarioID, perfilCodigoConta)
		}
		return perfilService.PerfilUsuario{}, err
	}
	var u perfilService.PerfilUsuario
	u.Identificador = usuarioID
	u.ContextoPerfil = "organization"
	u.URLImagemAvatar = textoEmMapa(detalhes, "avatar_image_url")
	u.URLImagemCapa = textoEmMapa(detalhes, "cover_image_url")
	u.Cargo = ""
	u.Curso = ""
	u.Semestre = ""
	u.TotalCandidaturas = 0
	u.TotalGrupos = 0
	u.TotalEventos = 0
	u.DestaqueComunidade = nil

	u.Email = textoEmMapa(detalhes, "email")
	u.CidadeEstado = juntarCidadeEstado(textoEmMapa(detalhes, "city"), textoEmMapa(detalhes, "state"))

	switch perfilCodigoConta {
	case "comunidade":
		if _, nome, kind, descricao, ok, err := repositorio.primeiraComunidadeDoUsuario(contexto, usuarioID); err != nil {
			return perfilService.PerfilUsuario{}, err
		} else if ok {
			u.Nome = nome
			u.SobreMim = descricao
			u.Instituicao = kind
		} else {
			u.Nome = textoEmMapa(detalhes, "community_name")
			u.Instituicao = textoEmMapa(detalhes, "community_type")
			u.SobreMim = ""
		}
	case "empresa":
		u.Nome = textoEmMapa(detalhes, "company_name")
		u.SobreMim = textoEmMapa(detalhes, "company_description")
		u.Instituicao = textoEmMapa(detalhes, "company_cnpj")
	case "universidade":
		u.Nome = textoEmMapa(detalhes, "institution_name")
		u.SobreMim = textoEmMapa(detalhes, "institution_description")
		sigla := textoEmMapa(detalhes, "institution_acronym")
		tipo := textoEmMapa(detalhes, "institution_type")
		u.Instituicao = strings.TrimSpace(strings.Join(nonEmptyParts(sigla, tipo), " · "))
	}

	u.Interesses = stringsFromDetalhesLista(detalhes, "interests")
	u.TopicosFavoritos = stringsFromDetalhesLista(detalhes, "favorite_topics")
	u.Especialidades = stringsFromDetalhesLista(detalhes, "specialties")
	if u.Interesses == nil {
		u.Interesses = []string{}
	}
	if u.TopicosFavoritos == nil {
		u.TopicosFavoritos = []string{}
	}
	if u.Especialidades == nil {
		u.Especialidades = []string{}
	}
	if err := repositorio.enriquecerPainelOrganizacao(contexto, &u, usuarioID, perfilCodigoConta, detalhes); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	repositorio.aplicarContratoFrontPerfil(contexto, &u, usuarioID, perfilCodigoConta)
	return u, nil
}

func (repositorio *perfilRepositoryPostgres) AtualizarPerfilParaExibicao(contexto context.Context, usuarioID, perfilCodigoConta string, corpo perfilService.RequisicaoAtualizarPerfil) (perfilService.PerfilUsuario, error) {
	switch perfilCodigoConta {
	case "comunidade", "empresa", "universidade":
		if err := repositorio.atualizarPerfilInstitucional(contexto, usuarioID, perfilCodigoConta, corpo); err != nil {
			return perfilService.PerfilUsuario{}, err
		}
		return repositorio.perfilInstitucional(contexto, usuarioID, perfilCodigoConta)
	default:
		return repositorio.AtualizarPerfilUsuario(contexto, usuarioID, corpo)
	}
}

func (repositorio *perfilRepositoryPostgres) atualizarPerfilInstitucional(contexto context.Context, usuarioID, perfilCodigoConta string, corpo perfilService.RequisicaoAtualizarPerfil) error {
	detalhes, err := repositorio.carregarDetalhesCadastro(contexto, usuarioID)
	if err != nil {
		return err
	}
	interesses := normalizarTags(corpo.Interesses)
	topicos := normalizarTags(corpo.TopicosFavoritos)
	especialidades := normalizarTags(corpo.Especialidades)
	detalhes["interests"] = interesses
	detalhes["favorite_topics"] = topicos
	detalhes["specialties"] = especialidades
	if corpo.URLImagemAvatar != nil {
		detalhes["avatar_image_url"] = strings.TrimSpace(*corpo.URLImagemAvatar)
	}
	if corpo.URLImagemCapa != nil {
		detalhes["cover_image_url"] = strings.TrimSpace(*corpo.URLImagemCapa)
	}

	switch perfilCodigoConta {
	case "comunidade":
		if v := strings.TrimSpace(corpo.Instituicao); v != "" {
			detalhes["community_name"] = v
		}
		if v := strings.TrimSpace(corpo.Curso); v != "" {
			detalhes["community_type"] = v
		}
		if _, _, _, _, ok, err := repositorio.primeiraComunidadeDoUsuario(contexto, usuarioID); err != nil {
			return err
		} else if ok {
			nome := strings.TrimSpace(corpo.Instituicao)
			if nome == "" {
				nome = textoEmMapa(detalhes, "community_name")
			}
			if err := repositorio.atualizarPrimeiraComunidadeDoUsuario(contexto, usuarioID, nome, strings.TrimSpace(corpo.SobreMim)); err != nil {
				return err
			}
		}
	case "empresa":
		if strings.TrimSpace(corpo.Instituicao) != "" {
			detalhes["company_name"] = strings.TrimSpace(corpo.Instituicao)
		}
		detalhes["company_description"] = strings.TrimSpace(corpo.SobreMim)
		if strings.TrimSpace(corpo.Cargo) != "" {
			detalhes["company_cnpj"] = strings.TrimSpace(corpo.Cargo)
		}
	case "universidade":
		if strings.TrimSpace(corpo.Instituicao) != "" {
			detalhes["institution_name"] = strings.TrimSpace(corpo.Instituicao)
		}
		detalhes["institution_description"] = strings.TrimSpace(corpo.SobreMim)
		if strings.TrimSpace(corpo.Curso) != "" {
			detalhes["institution_acronym"] = strings.TrimSpace(corpo.Curso)
		}
		if strings.TrimSpace(corpo.Semestre) != "" {
			detalhes["institution_type"] = strings.TrimSpace(corpo.Semestre)
		}
		if strings.TrimSpace(corpo.MapURL) != "" {
			detalhes["map_url"] = strings.TrimSpace(corpo.MapURL)
		}
	}

	js, err := json.Marshal(detalhes)
	if err != nil {
		return err
	}
	_, err = repositorio.pool.Exec(contexto, `
UPDATE cadastros_usuario SET details_json=$2::jsonb WHERE usuario_id=$1::uuid
`, usuarioID, js)
	return err
}

func (repositorio *perfilRepositoryPostgres) carregarDetalhesCadastro(contexto context.Context, usuarioID string) (map[string]any, error) {
	const sql = `SELECT details_json FROM cadastros_usuario WHERE usuario_id=$1::uuid`
	var raw []byte
	if err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&raw); err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if m == nil {
		m = map[string]any{}
	}
	return m, nil
}

func (repositorio *perfilRepositoryPostgres) primeiraComunidadeDoUsuario(contexto context.Context, usuarioID string) (id, nome, kind, descricao string, ok bool, err error) {
	const sql = `SELECT id::text, nome, kind, coalesce(description,'') FROM comunidades WHERE criado_por=$1::uuid ORDER BY criado_em ASC LIMIT 1`
	err = repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&id, &nome, &kind, &descricao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", "", "", false, nil
		}
		return "", "", "", "", false, err
	}
	return id, nome, kind, descricao, true, nil
}

func (repositorio *perfilRepositoryPostgres) perfilInstitucionalSemCadastro(contexto context.Context, usuarioID, perfilCodigoConta string) (perfilService.PerfilUsuario, error) {
	const sql = `SELECT coalesce(email,''), coalesce(avatar_image_url,''), coalesce(cover_image_url,'') FROM usuarios WHERE id=$1::uuid`
	var email, av, cv string
	if err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&email, &av, &cv); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	var u perfilService.PerfilUsuario
	u.Identificador = usuarioID
	u.ContextoPerfil = "organization"
	u.Email = email
	u.URLImagemAvatar = av
	u.URLImagemCapa = cv
	u.Interesses = []string{}
	u.TopicosFavoritos = []string{}
	u.Especialidades = []string{}
	detalhes := map[string]any{}
	if err := repositorio.enriquecerPainelOrganizacao(contexto, &u, usuarioID, perfilCodigoConta, detalhes); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	repositorio.aplicarContratoFrontPerfil(contexto, &u, usuarioID, perfilCodigoConta)
	return u, nil
}

func (repositorio *perfilRepositoryPostgres) atualizarPrimeiraComunidadeDoUsuario(contexto context.Context, usuarioID, nome, descricao string) error {
	const sql = `
UPDATE comunidades SET nome=$2, description=$3, atualizado_em=now()
WHERE id=(SELECT id FROM comunidades WHERE criado_por=$1::uuid ORDER BY criado_em ASC LIMIT 1)`
	ct, err := repositorio.pool.Exec(contexto, sql, usuarioID, nome, descricao)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return nil
}

func (repositorio *perfilRepositoryPostgres) AtualizarPerfilUsuario(contexto context.Context, usuarioID string, corpo perfilService.RequisicaoAtualizarPerfil) (perfilService.PerfilUsuario, error) {
	var avatar, cover string
	if err := repositorio.pool.QueryRow(contexto, `
SELECT coalesce(avatar_image_url,''), coalesce(cover_image_url,'') FROM usuarios WHERE id=$1::uuid`, usuarioID).Scan(&avatar, &cover); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	if corpo.URLImagemAvatar != nil {
		avatar = strings.TrimSpace(*corpo.URLImagemAvatar)
	}
	if corpo.URLImagemCapa != nil {
		cover = strings.TrimSpace(*corpo.URLImagemCapa)
	}
	interesses := normalizarTags(corpo.Interesses)
	topicos := normalizarTags(corpo.TopicosFavoritos)
	especialidades := normalizarTags(corpo.Especialidades)
	interessesJSON, _ := json.Marshal(interesses)
	topicosJSON, _ := json.Marshal(topicos)
	especialidadesJSON, _ := json.Marshal(especialidades)
	const sql = `
UPDATE usuarios SET
  about_me=$2,
  job_title=$3,
  course=$4,
  semester=$5,
  institution_name=$6,
  interests=$7::jsonb,
  favorite_topics=$8::jsonb,
  specialties=$9::jsonb,
  avatar_image_url=nullif($10,''),
  cover_image_url=nullif($11,''),
  course_and_semester=trim(both ' ' from concat_ws(' ', nullif($4,''), CASE WHEN nullif($5,'') IS NULL THEN '' ELSE concat('(', $5, 'º)') END)),
  atualizado_em=now()
WHERE id=$1::uuid`
	if _, err := repositorio.pool.Exec(contexto, sql, usuarioID, strings.TrimSpace(corpo.SobreMim), strings.TrimSpace(corpo.Cargo), strings.TrimSpace(corpo.Curso), strings.TrimSpace(corpo.Semestre), strings.TrimSpace(corpo.Instituicao), interessesJSON, topicosJSON, especialidadesJSON, avatar, cover); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	return repositorio.PerfilUsuario(contexto, usuarioID)
}

func (repositorio *perfilRepositoryPostgres) AtualizarURLImagemPerfil(contexto context.Context, usuarioID, perfilCodigoConta, tipoImagem, urlPublica string) error {
	coluna := "avatar_image_url"
	if tipoImagem == "cover" {
		coluna = "cover_image_url"
	}
	switch perfilCodigoConta {
	case "comunidade", "empresa", "universidade":
		const updCadastro = `
UPDATE cadastros_usuario
SET details_json = jsonb_set(COALESCE(details_json, '{}'::jsonb), $2::text[], to_jsonb($3::text), true)
WHERE usuario_id = $1::uuid`
		caminho := []string{coluna}
		ct, err := repositorio.pool.Exec(contexto, updCadastro, usuarioID, caminho, urlPublica)
		if err != nil {
			return err
		}
		if ct.RowsAffected() > 0 {
			return nil
		}
	}
	const updUsuario = `
UPDATE usuarios SET avatar_image_url = CASE WHEN $2 = 'avatar' THEN $3 ELSE avatar_image_url END,
                    cover_image_url = CASE WHEN $2 = 'cover' THEN $3 ELSE cover_image_url END,
                    atualizado_em = now()
WHERE id = $1::uuid`
	_, err := repositorio.pool.Exec(contexto, updUsuario, usuarioID, tipoImagem, urlPublica)
	return err
}

func (repositorio *perfilRepositoryPostgres) HistoricoPerfilUsuario(contexto context.Context, usuarioID string, limite int) (perfilService.RespostaHistoricoPerfil, error) {
	const sql = `
SELECT id, kind, title, subtitle, reference_id, created_at
FROM (
  SELECT p.id::text AS id, 'post'::text AS kind,
         left(p.body_text, 80) AS title, 'Post'::text AS subtitle,
         p.id::text AS reference_id, p.criado_em AS created_at
  FROM feed_posts p
  WHERE p.author_id=$1::uuid
  UNION ALL
  SELECT l.id::text, 'reading'::text, l.titulo, l.source, l.id::text, l.criado_em
  FROM leitura_semanal l
  WHERE l.criado_por=$1::uuid
  UNION ALL
  SELECT g.id::text, 'group'::text, g.titulo, g.field_of_study, g.id::text, g.criado_em
  FROM grupos_estudo g
  WHERE g.criado_por=$1::uuid
) h
ORDER BY created_at DESC
LIMIT $2`
	rows, err := repositorio.pool.Query(contexto, sql, usuarioID, limite)
	if err != nil {
		return perfilService.RespostaHistoricoPerfil{}, err
	}
	defer rows.Close()
	out := perfilService.RespostaHistoricoPerfil{}
	for rows.Next() {
		var it perfilService.ItemHistoricoPerfil
		var criadoEm time.Time
		if err := rows.Scan(&it.Identificador, &it.Tipo, &it.Titulo, &it.Subtitulo, &it.IDReferencia, &criadoEm); err != nil {
			return perfilService.RespostaHistoricoPerfil{}, err
		}
		it.CriadoEm = criadoEm.UTC().Format(time.RFC3339)
		out.Itens = append(out.Itens, it)
	}
	return out, rows.Err()
}

func (repositorio *perfilRepositoryPostgres) obterDestaqueComunidade(contexto context.Context, usuarioID string) (*perfilService.DestaqueComunidadePerfil, error) {
	const sql = `
SELECT c.id::text, c.nome, c.kind, uc.role
FROM usuario_comunidades uc
JOIN comunidades c ON c.id = uc.comunidade_id
WHERE uc.usuario_id=$1::uuid
ORDER BY uc.criado_em DESC
LIMIT 1`
	var destaque perfilService.DestaqueComunidadePerfil
	err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&destaque.Identificador, &destaque.Nome, &destaque.Tipo, &destaque.Papel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &destaque, nil
}

func normalizarTags(tags []string) []string {
	vistos := map[string]struct{}{}
	var out []string
	for _, t := range tags {
		tag := strings.TrimSpace(t)
		if tag == "" {
			continue
		}
		chave := strings.ToLower(tag)
		if _, ok := vistos[chave]; ok {
			continue
		}
		vistos[chave] = struct{}{}
		out = append(out, tag)
		if len(out) >= 20 {
			break
		}
	}
	return out
}

func textoEmMapa(m map[string]any, chave string) string {
	v, ok := m[chave]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strings.TrimSpace(fmt.Sprintf("%.0f", t))
	case json.Number:
		return strings.TrimSpace(t.String())
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func stringsFromDetalhesLista(m map[string]any, chave string) []string {
	v, ok := m[chave]
	if !ok || v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var out []string
	if err := json.Unmarshal(b, &out); err != nil {
		return nil
	}
	return out
}

func juntarCidadeEstado(cidade, estado string) string {
	cidade = strings.TrimSpace(cidade)
	estado = strings.TrimSpace(estado)
	if cidade == "" && estado == "" {
		return ""
	}
	if cidade == "" {
		return estado
	}
	if estado == "" {
		return cidade
	}
	return cidade + " - " + estado
}

func nonEmptyParts(parts ...string) []string {
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
