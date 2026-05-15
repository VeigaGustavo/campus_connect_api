package service

import "errors"

var (
	ErrPostInvalido       = errors.New("post invalido")
	ErrAnexoInvalido      = errors.New("anexo invalido")
	ErrCampoArquivoAusente = errors.New("campo file ausente")
)
