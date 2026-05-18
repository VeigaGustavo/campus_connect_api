package handler

import (
	"encoding/json"
	"net/http"

	eventoService "campus_connect_api/internal/modulos/evento/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type EventoHTTPHandler struct {
	servicoEvento *eventoService.EventoService
}

func NovoEventoHTTPHandler(servicoEvento *eventoService.EventoService) *EventoHTTPHandler {
	return &EventoHTTPHandler{servicoEvento: servicoEvento}
}

func (handler *EventoHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/events", respostas.AdaptadorHTTP(handler.GETEventos))
	grupo.POST("/events", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "sistema_admin")(handler.POSTCriarEvento)))
	grupo.PUT("/events/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "sistema_admin")(handler.PUTEvento)))
	grupo.DELETE("/events/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("universidade", "comunidade", "sistema_admin")(handler.DELETEEvento)))
}

func (handler *EventoHTTPHandler) GETEventos(resposta http.ResponseWriter, requisicao *http.Request) {
	eventos, err := handler.servicoEvento.ListarEventos(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, eventos)
}

func (handler *EventoHTTPHandler) POSTCriarEvento(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo eventoService.RequisicaoEvento
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	eventoCriado, err := handler.servicoEvento.CriarEvento(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, eventoCriado)
}

func (handler *EventoHTTPHandler) PUTEvento(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo eventoService.RequisicaoEvento
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	eventoAtualizado, err := handler.servicoEvento.AtualizarEvento(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, eventoAtualizado)
}

func (handler *EventoHTTPHandler) DELETEEvento(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoEvento.RemoverEvento(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

