package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	feedMedia "campus_connect_api/internal/modulos/feed/media"
	feedService "campus_connect_api/internal/modulos/feed/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
)

const maxMultipartFeedAnexo = 84 << 20 // 80 MiB video + margem form

func (handler *FeedHTTPHandler) POSTAnexoFeed(resposta http.ResponseWriter, requisicao *http.Request) {
	if _, ok := auth.SessaoDaRequisicao(requisicao); !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	if err := requisicao.ParseMultipartForm(maxMultipartFeedAnexo); err != nil {
		respostas.EscreverErro(resposta, http.StatusRequestEntityTooLarge, "file_too_large", "file too large")
		return
	}
	arquivo, nomeOriginal, err := primeiroArquivoFeed(requisicao)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_file", err.Error())
		return
	}
	defer arquivo.Close()

	tipo := strings.ToLower(strings.TrimSpace(requisicao.FormValue("type")))
	out, err := handler.servicoFeed.EnviarAnexoFeed(requisicao.Context(), tipo, nomeOriginal, arquivo)
	if err != nil {
		status, codigo, msg := mapearErroAnexoFeed(err)
		respostas.EscreverErro(resposta, status, codigo, msg)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func primeiroArquivoFeed(requisicao *http.Request) (multipart.File, string, error) {
	f, header, err := requisicao.FormFile("file")
	if err == nil && f != nil {
		nome := ""
		if header != nil {
			nome = header.Filename
		}
		return f, nome, nil
	}
	return nil, "", feedService.ErrCampoArquivoAusente
}

func mapearErroAnexoFeed(err error) (status int, codigo, msg string) {
	switch {
	case errors.Is(err, feedMedia.ErrArquivoGrande):
		return http.StatusRequestEntityTooLarge, "file_too_large", err.Error()
	case errors.Is(err, feedMedia.ErrFormatoInvalido), errors.Is(err, feedMedia.ErrArquivoVazio):
		return http.StatusBadRequest, "invalid_file", err.Error()
	case errors.Is(err, feedMedia.ErrTipoAnexo), errors.Is(err, feedService.ErrCampoArquivoAusente):
		return http.StatusBadRequest, "invalid_file", err.Error()
	default:
		return http.StatusInternalServerError, "server_error", err.Error()
	}
}
