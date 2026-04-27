package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	comum "campus_connect_api/internal/modulos/comum"
	projetoService "campus_connect_api/internal/modulos/projeto/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type ProjetoHTTPHandler struct {
	servicoProjeto *projetoService.ProjetoService
}

func NovoProjetoHTTPHandler(servicoProjeto *projetoService.ProjetoService) *ProjetoHTTPHandler {
	return &ProjetoHTTPHandler{servicoProjeto: servicoProjeto}
}

func (handler *ProjetoHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/projects", respostas.AdaptadorHTTP(handler.GETProjetos))
	grupo.POST("/projects", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTCriarProjeto)))
	grupo.PUT("/projects/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.PUTProjeto)))
	grupo.DELETE("/projects/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.DELETEProjeto)))
}

func (handler *ProjetoHTTPHandler) GETProjetos(resposta http.ResponseWriter, requisicao *http.Request) {
	out, err := handler.servicoProjeto.ListarProjetos(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *ProjetoHTTPHandler) POSTCriarProjeto(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo projetoService.RequisicaoProjeto
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	pr, err := handler.servicoProjeto.CriarProjeto(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, pr)
}

func (handler *ProjetoHTTPHandler) PUTProjeto(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo projetoService.RequisicaoProjeto
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	pr, err := handler.servicoProjeto.AtualizarProjeto(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, pr)
}

func (handler *ProjetoHTTPHandler) DELETEProjeto(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoProjeto.RemoverProjeto(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *ProjetoHTTPHandler) escreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, comum.ErrNaoEncontrado):
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, comum.ErrProibido):
		respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
