package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	comum "campus_connect_api/internal/modulos/comum"
	leituraService "campus_connect_api/internal/modulos/leitura/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type LeituraHTTPHandler struct {
	servicoLeitura *leituraService.LeituraService
}

func NovoLeituraHTTPHandler(servicoLeitura *leituraService.LeituraService) *LeituraHTTPHandler {
	return &LeituraHTTPHandler{servicoLeitura: servicoLeitura}
}

func (handler *LeituraHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/reading/weekly", respostas.AdaptadorHTTP(handler.GETLeituraSemanal))
	grupo.POST("/reading/weekly", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "empresa", "sistema_admin")(handler.POSTCriarLeituraSemanal)))
	grupo.PUT("/reading/weekly/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "empresa", "sistema_admin")(handler.PUTLeituraSemanal)))
	grupo.DELETE("/reading/weekly/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "empresa", "sistema_admin")(handler.DELETELeituraSemanal)))
}

func (handler *LeituraHTTPHandler) GETLeituraSemanal(resposta http.ResponseWriter, requisicao *http.Request) {
	out, err := handler.servicoLeitura.ListarLeituraSemanal(requisicao.Context(), requisicao.URL.Query().Get("kind"))
	if err != nil {
		if errors.Is(err, leituraService.ErrLeituraInvalida) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_reading", err.Error())
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if out == nil {
		out = []leituraService.ItemLeituraSemanal{}
	}
	respostas.EscreverJSON(resposta, http.StatusOK, leituraService.RespostaListaLeituraSemanal{Itens: out})
}

func (handler *LeituraHTTPHandler) POSTCriarLeituraSemanal(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo leituraService.RequisicaoLeituraSemanal
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	it, err := handler.servicoLeitura.CriarLeituraSemanal(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		if errors.Is(err, leituraService.ErrLeituraInvalida) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_reading", err.Error())
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, it)
}

func (handler *LeituraHTTPHandler) PUTLeituraSemanal(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo leituraService.RequisicaoLeituraSemanal
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	it, err := handler.servicoLeitura.AtualizarLeituraSemanal(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		if errors.Is(err, leituraService.ErrLeituraInvalida) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_reading", err.Error())
			return
		}
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, it)
}

func (handler *LeituraHTTPHandler) DELETELeituraSemanal(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoLeitura.RemoverLeituraSemanal(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		handler.escreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *LeituraHTTPHandler) escreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, comum.ErrNaoEncontrado):
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, comum.ErrProibido):
		respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
