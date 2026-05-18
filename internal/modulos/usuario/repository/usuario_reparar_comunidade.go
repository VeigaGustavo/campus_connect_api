package repository

import (
	"context"
	"encoding/json"
	"log"

	usuarioService "campus_connect_api/internal/modulos/usuario/service"
)

func (repositorio *usuarioRepositoryPostgres) RepararComunidadesSemGrupo(contexto context.Context) (int, error) {
	const sql = `
SELECT cu.usuario_id::text, cu.details_json
FROM cadastros_usuario cu
WHERE cu.profile_type = 'comunidade'
  AND NOT EXISTS (SELECT 1 FROM grupos_estudo g WHERE g.criado_por = cu.usuario_id)`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	reparados := 0
	for rows.Next() {
		var usuarioID string
		var detalhesJSON []byte
		if err := rows.Scan(&usuarioID, &detalhesJSON); err != nil {
			return reparados, err
		}
		var detalhes map[string]any
		if err := json.Unmarshal(detalhesJSON, &detalhes); err != nil {
			log.Printf("reparar comunidade: json usuario %s: %v", usuarioID, err)
			continue
		}
		req := usuarioService.RequisicaoCadastroUsuario{
			TipoPerfil:        "comunidade",
			TipoComunidade:    stringDe(detalhes, "community_type"),
			NomeComunidade:    stringDe(detalhes, "community_name"),
			TituloGrupo:       stringDe(detalhes, "group_title"),
			DescricaoGrupo:    stringDe(detalhes, "group_description"),
			VisibilidadeGrupo: stringDe(detalhes, "group_visibility"),
		}
		if req.NomeComunidade == "" {
			req.NomeComunidade = "Comunidade"
		}
		if req.DescricaoGrupo == "" {
			req.DescricaoGrupo = req.NomeComunidade
		}
		if req.VisibilidadeGrupo == "" {
			req.VisibilidadeGrupo = "public"
		}
		tx, err := repositorio.pool.Begin(contexto)
		if err != nil {
			return reparados, err
		}
		if _, err := repositorio.criarComunidadeEGrupoNoCadastro(contexto, tx, usuarioID, req); err != nil {
			_ = tx.Rollback(contexto)
			log.Printf("reparar comunidade: usuario %s: %v", usuarioID, err)
			continue
		}
		if err := tx.Commit(contexto); err != nil {
			return reparados, err
		}
		reparados++
	}
	return reparados, rows.Err()
}

func stringDe(m map[string]any, chave string) string {
	if m == nil {
		return ""
	}
	v, _ := m[chave].(string)
	return v
}
