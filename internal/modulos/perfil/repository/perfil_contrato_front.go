package repository

import (
	"context"
	"errors"
	"strings"

	"campus_connect_api/internal/modulos/perfil/media"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
	"github.com/jackc/pgx/v5"
)

func (repositorio *perfilRepositoryPostgres) obterPapelContaUsuario(contexto context.Context, usuarioID string) string {
	const sql = `SELECT pf.codigo FROM usuarios u JOIN perfis_usuario pf ON pf.id = u.perfil_id WHERE u.id=$1::uuid`
	var codigo string
	if err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&codigo); err != nil {
		return "padrao"
	}
	return codigo
}

func (repositorio *perfilRepositoryPostgres) aplicarContratoFrontPerfil(contexto context.Context, u *perfilService.PerfilUsuario, usuarioID, perfilCodigoConta string) {
	u.PapelConta = perfilCodigoConta
	u.TipoPerfil = repositorio.resolverTipoPerfil(contexto, usuarioID, perfilCodigoConta)
	u.URLImagemAvatar = media.AbsolutizarURL(u.URLImagemAvatar)
	u.URLImagemCapa = media.AbsolutizarURL(u.URLImagemCapa)
}

func (repositorio *perfilRepositoryPostgres) resolverTipoPerfil(contexto context.Context, usuarioID, perfilCodigoConta string) string {
	const sql = `SELECT profile_type FROM cadastros_usuario WHERE usuario_id=$1::uuid`
	var tipo string
	if err := repositorio.pool.QueryRow(contexto, sql, usuarioID).Scan(&tipo); err == nil {
		tipo = strings.TrimSpace(tipo)
		if tipo != "" {
			return tipo
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return tipoPerfilPadraoPorPapel(perfilCodigoConta)
	}
	return tipoPerfilPadraoPorPapel(perfilCodigoConta)
}

func tipoPerfilPadraoPorPapel(perfilCodigo string) string {
	switch perfilCodigo {
	case "comunidade":
		return "comunidade"
	case "empresa":
		return "empresa"
	case "universidade":
		return "universidade"
	default:
		return "estudante"
	}
}
