package database

import "errors"

var ErrNaoEncontrado = errors.New("nao encontrado")
var ErrProibido = errors.New("proibido")

// UsuarioInterno dados mínimos após autenticação ou criação.
type UsuarioInterno struct {
	ID           string
	Nome         string
	Email        string
	PerfilCodigo string
}

// Contrato completo de persistência: ver contratos.go (interface Armazenamento composta por ports).
