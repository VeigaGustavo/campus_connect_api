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

func ObrigarAutenticacao(proximo http.HandlerFunc) http.HandlerFunc {
	return func(resposta http.ResponseWriter, requisicao *http.Request) {
		cabecalhoAutorizacao := requisicao.Header.Get("Authorization")
		if !strings.HasPrefix(cabecalhoAutorizacao, "Bearer ") {
			respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(cabecalhoAutorizacao, "Bearer "))
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
