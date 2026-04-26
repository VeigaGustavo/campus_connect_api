package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"campus_connect_api/internal/infra/database"
	comunidadeService "campus_connect_api/internal/modulos/comunidade/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type ComunidadeHTTPHandler struct {
	servicoComunidade *comunidadeService.ComunidadeService
}

func NovoComunidadeHTTPHandler(servicoComunidade *comunidadeService.ComunidadeService) *ComunidadeHTTPHandler {
	return &ComunidadeHTTPHandler{servicoComunidade: servicoComunidade}
}

func (handler *ComunidadeHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/communities", respostas.AdaptadorHTTP(handler.GETComunidades))
	grupo.POST("/communities", respostas.AdaptadorHTTP(auth.ExigirPerfis("comunidade", "sistema_admin")(handler.POSTCriarComunidade)))
	grupo.PUT("/communities/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("comunidade", "sistema_admin")(handler.PUTComunidade)))
	grupo.DELETE("/communities/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("comunidade", "sistema_admin")(handler.DELETEComunidade)))
}

func (handler *ComunidadeHTTPHandler) GETComunidades(resposta http.ResponseWriter, requisicao *http.Request) {
	comunidades, err := handler.servicoComunidade.ListarComunidades(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, comunidades)
}

func (handler *ComunidadeHTTPHandler) POSTCriarComunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo comunidadeService.RequisicaoCriarComunidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	comunidadeCriada, err := handler.servicoComunidade.CriarComunidade(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, comunidadeCriada)
}

func (handler *ComunidadeHTTPHandler) PUTComunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo comunidadeService.RequisicaoCriarComunidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	comunidadeAtualizada, err := handler.servicoComunidade.AtualizarComunidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, comunidadeAtualizada)
}

func (handler *ComunidadeHTTPHandler) DELETEComunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoComunidade.RemoverComunidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *ComunidadeHTTPHandler) escreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNaoEncontrado):
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, database.ErrProibido):
		respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
