package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"campus_connect_api/internal/modulos/perfil/media"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

const maxMultipartPerfil = 6 << 20

func (handler *PerfilHTTPHandler) RegistrarRotasUploadGIN(grupo *gin.RouterGroup) {
	grupo.POST("/profile/avatar", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.POSTImagemPerfil("avatar"))))
	grupo.POST("/profile/cover", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.POSTImagemPerfil("cover"))))
}

func (handler *PerfilHTTPHandler) POSTImagemPerfil(tipo string) http.HandlerFunc {
	return func(resposta http.ResponseWriter, requisicao *http.Request) {
		sessao, ok := auth.SessaoDaRequisicao(requisicao)
		if !ok {
			respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
			return
		}
		if err := requisicao.ParseMultipartForm(maxMultipartPerfil); err != nil {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_multipart", "invalid multipart form")
			return
		}
		arquivo, _, err := primeiroArquivoMultipart(requisicao, "file", "image", "avatar", "cover")
		if err != nil {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "missing_file", err.Error())
			return
		}
		defer arquivo.Close()

		out, err := handler.servicoPerfil.EnviarImagemPerfil(requisicao.Context(), sessao.UsuarioID, sessao.Perfil, tipo, arquivo)
		if err != nil {
			status, codigo, msg := mapearErroUpload(err)
			respostas.EscreverErro(resposta, status, codigo, msg)
			return
		}
		respostas.EscreverJSON(resposta, http.StatusOK, out)
	}
}

func primeiroArquivoMultipart(requisicao *http.Request, nomes ...string) (multipart.File, string, error) {
	for _, nome := range nomes {
		f, header, err := requisicao.FormFile(nome)
		if err == nil && f != nil {
			return f, header.Filename, nil
		}
	}
	return nil, "", perfilService.ErrCampoArquivoAusente
}

func mapearErroUpload(err error) (status int, codigo, msg string) {
	switch {
	case errors.Is(err, media.ErrArquivoGrande):
		return http.StatusRequestEntityTooLarge, "file_too_large", err.Error()
	case errors.Is(err, media.ErrFormatoInvalido), errors.Is(err, media.ErrImagemVazia):
		return http.StatusUnsupportedMediaType, "invalid_image", err.Error()
	case errors.Is(err, perfilService.ErrCampoArquivoAusente):
		return http.StatusBadRequest, "missing_file", err.Error()
	case errors.Is(err, perfilService.ErrTipoImagemInvalido):
		return http.StatusBadRequest, "invalid_type", err.Error()
	default:
		if strings.Contains(err.Error(), "base64") {
			return http.StatusBadRequest, "invalid_image_url", err.Error()
		}
		return http.StatusInternalServerError, "server_error", err.Error()
	}
}
