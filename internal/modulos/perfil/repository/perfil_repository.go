package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

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
	return u, nil
}

func (repositorio *perfilRepositoryPostgres) AtualizarPerfilUsuario(contexto context.Context, usuarioID string, corpo perfilService.RequisicaoAtualizarPerfil) (perfilService.PerfilUsuario, error) {
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
  course_and_semester=trim(both ' ' from concat_ws(' ', nullif($4,''), CASE WHEN nullif($5,'') IS NULL THEN '' ELSE concat('(', $5, 'º)') END)),
  atualizado_em=now()
WHERE id=$1::uuid`
	if _, err := repositorio.pool.Exec(contexto, sql, usuarioID, strings.TrimSpace(corpo.SobreMim), strings.TrimSpace(corpo.Cargo), strings.TrimSpace(corpo.Curso), strings.TrimSpace(corpo.Semestre), strings.TrimSpace(corpo.Instituicao), interessesJSON, topicosJSON, especialidadesJSON); err != nil {
		return perfilService.PerfilUsuario{}, err
	}
	return repositorio.PerfilUsuario(contexto, usuarioID)
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
