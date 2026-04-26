package handler

import (
	"net/http"

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

func (handler *PerfilHTTPHandler) RegistrarRotasHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /profile", auth.ObrigarAutenticacao(handler.GETPerfil))
}

func (handler *PerfilHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/profile", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETPerfil)))
}

func (handler *PerfilHTTPHandler) GETPerfil(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	out, err := handler.servicoPerfil.ObterPerfil(requisicao.Context(), sessao.UsuarioID)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}
