package database

import "campus_connect_api/internal/modulos/comum"

var ErrNaoEncontrado = comum.ErrNaoEncontrado
var ErrProibido = comum.ErrProibido

// UsuarioInterno dados mínimos após autenticação ou criação.
type UsuarioInterno struct {
	ID           string
	Nome         string
	Email        string
	PerfilCodigo string
}

// Erros e tipos comuns usados pelo adapter de banco.
