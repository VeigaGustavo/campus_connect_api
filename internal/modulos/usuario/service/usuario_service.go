package service

import (
	"context"
	"errors"
	"fmt"
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
	corpo.TipoPerfil = strings.TrimSpace(corpo.TipoPerfil)
	corpo.TipoPerfil = strings.TrimPrefix(corpo.TipoPerfil, "\ufeff")
	corpo.TipoPerfil = strings.TrimSuffix(corpo.TipoPerfil, "\ufeff")
	corpo.TipoPerfil = strings.ToLower(corpo.TipoPerfil)
	corpo.DataNascimento = normalizarDataNascimento(corpo.DataNascimento)
	if corpo.Idade <= 0 {
		if idade, ok := idadeAPartirDeDataNascimento(strings.TrimSpace(corpo.DataNascimento)); ok {
			corpo.Idade = idade
		}
	}
	if strings.TrimSpace(corpo.NomeCompleto) == "" ||
		strings.TrimSpace(corpo.Email) == "" ||
		strings.TrimSpace(corpo.Senha) == "" ||
		strings.TrimSpace(corpo.CPF) == "" ||
		strings.TrimSpace(corpo.Cidade) == "" ||
		strings.TrimSpace(corpo.Estado) == "" ||
		corpo.Idade <= 0 {
		return nil, fmt.Errorf("%w: full_name,email,password,cpf,city,state e idade valida (age ou birth_date) sao obrigatorios; idade=%d birth_date=%q",
			ErrCadastroInvalido, corpo.Idade, corpo.DataNascimento)
	}

	switch corpo.TipoPerfil {
	case "estudante":
		if strings.TrimSpace(corpo.Instituicao) == "" {
			return nil, fmt.Errorf("%w: institution obrigatorio para estudante", ErrCadastroInvalido)
		}
	case "comunidade":
		if strings.TrimSpace(corpo.Instituicao) == "" {
			return nil, fmt.Errorf("%w: institution obrigatorio para comunidade", ErrCadastroInvalido)
		}
		if strings.TrimSpace(corpo.TipoComunidade) == "" || strings.TrimSpace(corpo.NomeComunidade) == "" {
			return nil, fmt.Errorf("%w: community_type e community_name obrigatorios", ErrCadastroInvalido)
		}
		if strings.TrimSpace(corpo.DescricaoGrupo) == "" {
			return nil, fmt.Errorf("%w: group_description obrigatorio", ErrCadastroInvalido)
		}
		vis := strings.TrimSpace(corpo.VisibilidadeGrupo)
		if vis != "public" && vis != "private" {
			return nil, fmt.Errorf("%w: group_visibility deve ser public ou private", ErrCadastroInvalido)
		}
	case "empresa":
		if strings.TrimSpace(corpo.NomeEmpresa) == "" {
			return nil, fmt.Errorf("%w: company_name obrigatorio para empresa", ErrCadastroInvalido)
		}
	case "universidade":
		if strings.TrimSpace(corpo.NomeInstituicao) == "" {
			return nil, fmt.Errorf("%w: institution_name obrigatorio para universidade", ErrCadastroInvalido)
		}
	default:
		return nil, fmt.Errorf("%w: profile_type invalido (recebido: %q)", ErrCadastroInvalido, corpo.TipoPerfil)
	}

	usuario, comunidade, err := servico.repositorio.CriarUsuarioComCadastro(contexto, corpo)
	if err != nil {
		return nil, err
	}
	out := map[string]any{
		"id": usuario.ID, "name": usuario.Nome, "email": usuario.Email, "role": usuario.PerfilCodigo, "profile_type": corpo.TipoPerfil,
	}
	if comunidade != nil {
		if comunidade.CommunityID != "" {
			out["community_id"] = comunidade.CommunityID
		}
		if comunidade.GroupID != "" {
			out["group_id"] = comunidade.GroupID
		}
	}
	return out, nil
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

// normalizarDataNascimento aceita YYYY-MM-DD ou ISO completo (ex.: Flutter DateTime.toIso8601String).
func normalizarDataNascimento(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC().Format("2006-01-02")
	}
	if i := strings.IndexByte(s, 'T'); i > 0 {
		prefix := s[:i]
		if t, err := time.Parse("2006-01-02", prefix); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return s
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
