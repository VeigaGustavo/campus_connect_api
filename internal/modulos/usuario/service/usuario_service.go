package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"campus_connect_api/internal/modelos"
	segurancaAuth "campus_connect_api/internal/modulos/seguranca/auth"
)

var (
	ErrCadastroInvalido = errors.New("cadastro invalido")
	ErrSemPermissao     = errors.New("sem permissao")
)

type UsuarioService struct {
	repositorio UsuarioRepository
}

func NovoUsuarioService(repositorio UsuarioRepository) *UsuarioService {
	return &UsuarioService{repositorio: repositorio}
}

func (servico *UsuarioService) RegistrarNovoUsuario(contexto context.Context, corpo RequisicaoCadastroUsuario) (map[string]any, error) {
	if corpo.Idade <= 0 {
		if idade, ok := idadeAPartirDeDataNascimento(strings.TrimSpace(corpo.DataNascimento)); ok {
			corpo.Idade = idade
		}
	}
	if strings.TrimSpace(corpo.NomeCompleto) == "" ||
		strings.TrimSpace(corpo.Email) == "" ||
		strings.TrimSpace(corpo.Senha) == "" ||
		strings.TrimSpace(corpo.CPF) == "" ||
		strings.TrimSpace(corpo.Instituicao) == "" ||
		strings.TrimSpace(corpo.Cidade) == "" ||
		strings.TrimSpace(corpo.Estado) == "" ||
		corpo.Idade <= 0 {
		return nil, ErrCadastroInvalido
	}

	switch corpo.TipoPerfil {
	case "estudante":
	case "comunidade":
		if strings.TrimSpace(corpo.TipoComunidade) == "" || strings.TrimSpace(corpo.NomeComunidade) == "" {
			return nil, ErrCadastroInvalido
		}
	case "empresa":
		if strings.TrimSpace(corpo.NomeEmpresa) == "" {
			return nil, ErrCadastroInvalido
		}
	case "universidade":
		if strings.TrimSpace(corpo.NomeInstituicao) == "" {
			return nil, ErrCadastroInvalido
		}
	default:
		return nil, ErrCadastroInvalido
	}

	usuario, err := servico.repositorio.CriarUsuarioComCadastro(contexto, corpo)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id": usuario.ID, "name": usuario.Nome, "email": usuario.Email, "role": usuario.PerfilCodigo, "profile_type": corpo.TipoPerfil,
	}, nil
}

func (servico *UsuarioService) CriarUsuarioAdministracao(contexto context.Context, sessao segurancaAuth.SessaoUsuario, corpo RequisicaoCriarUsuario) (*UsuarioInterno, error) {
	if sessao.Perfil != "sistema_admin" {
		return nil, ErrSemPermissao
	}
	return servico.repositorio.CriarUsuario(contexto, corpo.Nome, corpo.Email, corpo.Senha, corpo.Perfil)
}

func SessaoParaRespostaUsuario(sessao segurancaAuth.SessaoUsuario) modelos.UsuarioSessao {
	return modelos.UsuarioSessao{
		ID:     sessao.UsuarioID,
		Email:  sessao.Email,
		Perfil: sessao.Perfil,
	}
}

// idadeAPartirDeDataNascimento espera data no formato YYYY-MM-DD (RFC 3339 só a parte da data).
func idadeAPartirDeDataNascimento(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	nascimento, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, false
	}
	hoje := time.Now().UTC()
	if nascimento.After(hoje) {
		return 0, false
	}
	anos := hoje.Year() - nascimento.Year()
	if hoje.Month() < nascimento.Month() || (hoje.Month() == nascimento.Month() && hoje.Day() < nascimento.Day()) {
		anos--
	}
	if anos < 1 || anos > 130 {
		return 0, false
	}
	return anos, true
}
