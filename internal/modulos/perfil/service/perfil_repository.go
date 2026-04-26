package service

import (
	"context"
)

type PerfilRepository interface {
	PerfilUsuario(contexto context.Context, usuarioID string) (PerfilUsuario, error)
}
