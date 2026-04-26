package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"campus_connect_api/internal/infra/database"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	universidadeService "campus_connect_api/internal/modulos/universidade/service"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type UniversidadeHTTPHandler struct {
	servicoUniversidade *universidadeService.UniversidadeService
}

func NovoUniversidadeHTTPHandler(servicoUniversidade *universidadeService.UniversidadeService) *UniversidadeHTTPHandler {
	return &UniversidadeHTTPHandler{servicoUniversidade: servicoUniversidade}
}

func (handler *UniversidadeHTTPHandler) RegistrarRotasHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /university/notices", handler.GETAvisosUniversidade)
	mux.HandleFunc("POST /university/notices", auth.ExigirPerfis("universidade", "sistema_admin")(handler.POSTCriarAvisoUniversidade))
	mux.HandleFunc("PUT /university/notices/{id}", auth.ExigirPerfis("universidade", "sistema_admin")(handler.PUTAvisoUniversidade))
	mux.HandleFunc("DELETE /university/notices/{id}", auth.ExigirPerfis("universidade", "sistema_admin")(handler.DELETEAvisoUniversidade))
}

func (handler *UniversidadeHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/university/notices", respostas.AdaptadorHTTP(handler.GETAvisosUniversidade))
	grupo.POST("/university/notices", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "sistema_admin")(handler.POSTCriarAvisoUniversidade)))
	grupo.PUT("/university/notices/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "sistema_admin")(handler.PUTAvisoUniversidade)))
	grupo.DELETE("/university/notices/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "sistema_admin")(handler.DELETEAvisoUniversidade)))
}

func (handler *UniversidadeHTTPHandler) GETAvisosUniversidade(resposta http.ResponseWriter, requisicao *http.Request) {
	out, err := handler.servicoUniversidade.ListarAvisosUniversidade(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *UniversidadeHTTPHandler) POSTCriarAvisoUniversidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo universidadeService.RequisicaoCriarAvisoUniversidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	avisoCriado, err := handler.servicoUniversidade.CriarAvisoUniversidade(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, avisoCriado)
}

func (handler *UniversidadeHTTPHandler) PUTAvisoUniversidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo universidadeService.RequisicaoCriarAvisoUniversidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	avisoAtualizado, err := handler.servicoUniversidade.AtualizarAvisoUniversidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, avisoAtualizado)
}

func (handler *UniversidadeHTTPHandler) DELETEAvisoUniversidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoUniversidade.RemoverAvisoUniversidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *UniversidadeHTTPHandler) escreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNaoEncontrado):
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, database.ErrProibido):
		respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
