package service

import (
	"context"
	"errors"
)

type SegurancaRepository interface {
	Autenticar(contexto context.Context, email, senha string) (*UsuarioAutenticado, error)
}

var ErrAutenticacaoInvalida = errors.New("autenticacao invalida")

type UsuarioAutenticado struct {
	ID           string
	Nome         string
	Email        string
	PerfilCodigo string
}
