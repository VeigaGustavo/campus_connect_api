package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	perfilService "campus_connect_api/internal/modulos/perfil/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type PerfilHTTPHandler struct {
	servicoPerfil *perfilService.PerfilService
}

func NovoPerfilHTTPHandler(servicoPerfil *perfilService.PerfilService) *PerfilHTTPHandler {
	return &PerfilHTTPHandler{servicoPerfil: servicoPerfil}
}

func (handler *PerfilHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/profile", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETPerfil)))
	grupo.PUT("/profile", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.PUTPerfil)))
	grupo.GET("/profile/history", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETHistoricoPerfil)))
}

func (handler *PerfilHTTPHandler) GETPerfil(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	out, err := handler.servicoPerfil.ObterPerfil(requisicao.Context(), sessao.UsuarioID, sessao.Perfil)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *PerfilHTTPHandler) PUTPerfil(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo perfilService.RequisicaoAtualizarPerfil
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	out, err := handler.servicoPerfil.AtualizarPerfil(requisicao.Context(), sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *PerfilHTTPHandler) GETHistoricoPerfil(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	limite := 20
	if v := requisicao.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limite = n
		}
	}
	out, err := handler.servicoPerfil.HistoricoPerfil(requisicao.Context(), sessao.UsuarioID, limite)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}
