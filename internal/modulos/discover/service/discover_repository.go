package service

import (
	"context"
)

type DiscoverRepository interface {
	FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]ItemDescobrir, error)
}
