package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	comum "campus_connect_api/internal/modulos/comum"
	empresaService "campus_connect_api/internal/modulos/empresa/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type EmpresaHTTPHandler struct {
	servicoEmpresa *empresaService.EmpresaService
}

func NovoEmpresaHTTPHandler(servicoEmpresa *empresaService.EmpresaService) *EmpresaHTTPHandler {
	return &EmpresaHTTPHandler{servicoEmpresa: servicoEmpresa}
}

func (handler *EmpresaHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/opportunities/:id", respostas.AdaptadorHTTP(handler.GETOportunidadePorID))
	grupo.GET("/opportunities/:id/applicants", respostas.AdaptadorHTTP(auth.ExigirPerfis("empresa", "sistema_admin")(handler.GETCandidatosOportunidade)))
	grupo.GET("/opportunities", respostas.AdaptadorHTTP(handler.GETOportunidades))
	grupo.POST("/opportunities", respostas.AdaptadorHTTP(auth.ExigirPerfis("empresa", "sistema_admin")(handler.POSTCriarOportunidade)))
	grupo.PUT("/opportunities/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("empresa", "sistema_admin")(handler.PUTOportunidade)))
	grupo.DELETE("/opportunities/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("empresa", "sistema_admin")(handler.DELETEOportunidade)))
}

func (handler *EmpresaHTTPHandler) GETOportunidades(resposta http.ResponseWriter, requisicao *http.Request) {
	oportunidades, err := handler.servicoEmpresa.ListarOportunidades(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, oportunidades)
}

func (handler *EmpresaHTTPHandler) GETOportunidadePorID(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	oportunidade, encontrada, err := handler.servicoEmpresa.ObterOportunidadePorID(requisicao.Context(), id)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if !encontrada {
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, oportunidade)
}

func (handler *EmpresaHTTPHandler) GETCandidatosOportunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	respostas.EscreverJSON(resposta, http.StatusOK, handler.servicoEmpresa.ListarCandidatosOportunidade(id))
}

func (handler *EmpresaHTTPHandler) POSTCriarOportunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo empresaService.RequisicaoCriarOportunidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	oportunidadeCriada, err := handler.servicoEmpresa.CriarOportunidade(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		if errors.Is(err, empresaService.ErrOportunidadeInvalida) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_opportunity", err.Error())
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, oportunidadeCriada)
}

func (handler *EmpresaHTTPHandler) PUTOportunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo empresaService.RequisicaoCriarOportunidade
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	oportunidadeAtualizada, err := handler.servicoEmpresa.AtualizarOportunidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		if errors.Is(err, empresaService.ErrOportunidadeInvalida) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_opportunity", err.Error())
			return
		}
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, oportunidadeAtualizada)
}

func (handler *EmpresaHTTPHandler) DELETEOportunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoEmpresa.RemoverOportunidade(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *EmpresaHTTPHandler) escreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, comum.ErrNaoEncontrado):
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, comum.ErrProibido):
		respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
