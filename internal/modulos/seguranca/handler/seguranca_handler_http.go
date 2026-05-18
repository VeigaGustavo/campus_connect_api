package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"campus_connect_api/internal/modelos"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	segurancaService "campus_connect_api/internal/modulos/seguranca/service"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type SegurancaHTTPHandler struct {
	servicoSeguranca *segurancaService.SegurancaService
}

func NovoSegurancaHTTPHandler(servicoSeguranca *segurancaService.SegurancaService) *SegurancaHTTPHandler {
	return &SegurancaHTTPHandler{servicoSeguranca: servicoSeguranca}
}

func (handler *SegurancaHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.POST("/auth/login", respostas.AdaptadorHTTP(handler.POSTLogin))
	grupo.GET("/auth/me", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETSessaoAtual)))
}

func (handler *SegurancaHTTPHandler) POSTLogin(resposta http.ResponseWriter, requisicao *http.Request) {
	var corpo modelos.RequisicaoLogin
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	out, err := handler.servicoSeguranca.RealizarLogin(requisicao.Context(), corpo.Email, corpo.Senha)
	if err != nil {
		if errors.Is(err, segurancaService.ErrCredenciaisInvalidas) {
			respostas.EscreverErro(resposta, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "token_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *SegurancaHTTPHandler) GETSessaoAtual(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, usuarioService.SessaoParaRespostaUsuario(sessao))
}
