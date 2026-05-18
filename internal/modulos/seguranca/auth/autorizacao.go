package auth

import (
	"context"
	"net/http"
	"strings"

	segurancaService "campus_connect_api/internal/modulos/seguranca/service"
	"campus_connect_api/internal/respostas"
)

type chaveContexto string

const chaveSessao chaveContexto = "sessao_usuario"

type SessaoUsuario struct {
	UsuarioID string
	Email     string
	Perfil    string
}

func extrairToken(requisicao *http.Request) string {
	cabecalhoAutorizacao := requisicao.Header.Get("Authorization")
	if strings.HasPrefix(cabecalhoAutorizacao, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(cabecalhoAutorizacao, "Bearer "))
	}
	if t := strings.TrimSpace(requisicao.URL.Query().Get("access_token")); t != "" {
		return t
	}
	return strings.TrimSpace(requisicao.URL.Query().Get("token"))
}

func ObrigarAutenticacao(proximo http.HandlerFunc) http.HandlerFunc {
	return func(resposta http.ResponseWriter, requisicao *http.Request) {
		token := extrairToken(requisicao)
		if token == "" {
			respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			return
		}
		payload, err := segurancaService.ValidarToken(token)
		if err != nil {
			respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "invalid token")
			return
		}
		sessao := SessaoUsuario{UsuarioID: payload.UsuarioID, Email: payload.Email, Perfil: payload.Perfil}
		contextoRequisicao := context.WithValue(requisicao.Context(), chaveSessao, sessao)
		proximo(resposta, requisicao.WithContext(contextoRequisicao))
	}
}

func ExigirPerfis(perfisPermitidos ...string) func(http.HandlerFunc) http.HandlerFunc {
	permitidos := map[string]struct{}{}
	for _, p := range perfisPermitidos {
		permitidos[p] = struct{}{}
	}
	return func(proximo http.HandlerFunc) http.HandlerFunc {
		return ObrigarAutenticacao(func(resposta http.ResponseWriter, requisicao *http.Request) {
			sessao, sessaoEncontrada := SessaoDaRequisicao(requisicao)
			if !sessaoEncontrada {
				respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "user session not found")
				return
			}
			if _, existe := permitidos[sessao.Perfil]; !existe {
				respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "insufficient role")
				return
			}
			proximo(resposta, requisicao)
		})
	}
}

func SessaoDaRequisicao(requisicao *http.Request) (SessaoUsuario, bool) {
	valorSessao := requisicao.Context().Value(chaveSessao)
	sessao, ok := valorSessao.(SessaoUsuario)
	return sessao, ok
}
