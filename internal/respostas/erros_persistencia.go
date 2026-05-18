package respostas

import (
	"errors"
	"net/http"

	comum "campus_connect_api/internal/modulos/comum"
)

func EscreverErroPersistencia(resposta http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, comum.ErrNaoEncontrado):
		EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, comum.ErrProibido):
		EscreverErro(resposta, http.StatusForbidden, "forbidden", "not allowed")
	default:
		EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
	}
}
