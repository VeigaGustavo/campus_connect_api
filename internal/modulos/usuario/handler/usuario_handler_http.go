package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	auth "campus_connect_api/internal/modulos/seguranca/auth"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type UsuarioHTTPHandler struct {
	servicoUsuario *usuarioService.UsuarioService
}

func NovoUsuarioHTTPHandler(servicoUsuario *usuarioService.UsuarioService) *UsuarioHTTPHandler {
	return &UsuarioHTTPHandler{servicoUsuario: servicoUsuario}
}

func (handler *UsuarioHTTPHandler) RegistrarRotasHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", handler.POSTCadastroUsuarioPublico)
	mux.HandleFunc("POST /admin/users", auth.ExigirPerfis("sistema_admin")(handler.POSTCriarUsuarioAdmin))
}

func (handler *UsuarioHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.POST("/auth/register", respostas.AdaptadorHTTP(handler.POSTCadastroUsuarioPublico))
	grupo.POST("/admin/users", respostas.AdaptadorHTTP(auth.ExigirPerfis("sistema_admin")(handler.POSTCriarUsuarioAdmin)))
}

func (handler *UsuarioHTTPHandler) POSTCadastroUsuarioPublico(resposta http.ResponseWriter, requisicao *http.Request) {
	var corpo usuarioService.RequisicaoCadastroUsuario
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	out, err := handler.servicoUsuario.RegistrarNovoUsuario(requisicao.Context(), corpo)
	if err != nil {
		if errors.Is(err, usuarioService.ErrCadastroInvalido) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_registration", "missing or invalid fields")
			return
		}
		respostas.EscreverErro(resposta, http.StatusBadRequest, "registration_failed", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, out)
}

func (handler *UsuarioHTTPHandler) POSTCriarUsuarioAdmin(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo usuarioService.RequisicaoCriarUsuario
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	usuarioCriado, err := handler.servicoUsuario.CriarUsuarioAdministracao(requisicao.Context(), sessao, corpo)
	if err != nil {
		if errors.Is(err, usuarioService.ErrSemPermissao) {
			respostas.EscreverErro(resposta, http.StatusForbidden, "forbidden", "admin only")
			return
		}
		respostas.EscreverErro(resposta, http.StatusBadRequest, "create_user_failed", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, map[string]any{
		"id": usuarioCriado.ID, "email": usuarioCriado.Email, "role": usuarioCriado.PerfilCodigo, "name": usuarioCriado.Nome,
	})
}
