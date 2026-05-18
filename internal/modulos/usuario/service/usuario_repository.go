package service

import (
	"context"
)

type UsuarioRepository interface {
	CriarUsuario(contexto context.Context, nome, email, senha, perfilCodigo string) (*UsuarioInterno, error)
	CriarUsuarioComCadastro(contexto context.Context, requisicao RequisicaoCadastroUsuario) (*UsuarioInterno, *ResultadoCadastroComunidade, error)
	RepararComunidadesSemGrupo(contexto context.Context) (int, error)
}

type UsuarioInterno struct {
	ID           string
	Nome         string
	Email        string
	PerfilCodigo string
}

// ResultadoCadastroComunidade IDs criados no register de profile_type=comunidade.
type ResultadoCadastroComunidade struct {
	CommunityID string
	GroupID     string
}
