package repository

import (
	"context"
	"encoding/json"

	perfilService "campus_connect_api/internal/modulos/perfil/service"
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
SELECT nome, coalesce(initials,''), coalesce(cover_image_url,''), coalesce(avatar_image_url,''),
       coalesce(performance_certificate_label,''), coalesce(course_and_semester,''), email, coalesce(city_state,''),
       applications_count, groups_count, events_count,
       coalesce(interests,'[]'::jsonb), coalesce(recent_activity,'[]'::jsonb)
FROM usuarios WHERE id=$1::uuid`
	var u perfilService.PerfilUsuario
	var interestsJSON, recentJSON []byte
	err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(
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
