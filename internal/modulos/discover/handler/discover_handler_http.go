package handler

import (
	"net/http"
	"strings"

	discoverService "campus_connect_api/internal/modulos/discover/service"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type DiscoverHTTPHandler struct {
	servicoDiscover *discoverService.DiscoverService
}

func NovoDiscoverHTTPHandler(servicoDiscover *discoverService.DiscoverService) *DiscoverHTTPHandler {
	return &DiscoverHTTPHandler{servicoDiscover: servicoDiscover}
}

func (handler *DiscoverHTTPHandler) RegistrarRotasHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /discover", handler.GETDescobrir)
}

func (handler *DiscoverHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/discover", respostas.AdaptadorHTTP(handler.GETDescobrir))
}

func (handler *DiscoverHTTPHandler) GETDescobrir(resposta http.ResponseWriter, requisicao *http.Request) {
	filtro := requisicao.URL.Query().Get("filter")
	gruposCSV := strings.TrimSpace(requisicao.URL.Query().Get("group_ids"))
	var grupos []string
	if gruposCSV != "" {
		for _, g := range strings.Split(gruposCSV, ",") {
			gg := strings.TrimSpace(g)
			if gg != "" {
				grupos = append(grupos, gg)
			}
		}
	}
	out, err := handler.servicoDiscover.FeedDescobrir(requisicao.Context(), filtro, grupos)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}
