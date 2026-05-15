package service

import (
	"context"
	"errors"
	"io"
	"strings"

	"campus_connect_api/internal/modulos/perfil/media"
)

type PerfilService struct {
	repositorio PerfilRepository
}

func NovoPerfilService(repositorio PerfilRepository) *PerfilService {
	return &PerfilService{repositorio: repositorio}
}

func (servico *PerfilService) ObterPerfil(contexto context.Context, usuarioID, perfilCodigoConta string) (PerfilUsuario, error) {
	return servico.repositorio.PerfilParaExibicao(contexto, usuarioID, perfilCodigoConta)
}

func (servico *PerfilService) AtualizarPerfil(contexto context.Context, usuarioID, perfilCodigoConta string, corpo RequisicaoAtualizarPerfil) (PerfilUsuario, error) {
	if corpo.URLImagemAvatar != nil {
		if err := validarURLImagemOpcional(*corpo.URLImagemAvatar); err != nil {
			return PerfilUsuario{}, err
		}
	}
	if corpo.URLImagemCapa != nil {
		if err := validarURLImagemOpcional(*corpo.URLImagemCapa); err != nil {
			return PerfilUsuario{}, err
		}
	}
	return servico.repositorio.AtualizarPerfilParaExibicao(contexto, usuarioID, perfilCodigoConta, corpo)
}

func validarURLImagemOpcional(url string) error {
	u := strings.TrimSpace(url)
	if u == "" {
		return nil
	}
	if strings.HasPrefix(strings.ToLower(u), "data:") {
		return errors.New("nao envie imagem em base64 no PUT; use POST /api/profile/avatar ou /api/profile/cover")
	}
	if len(u) > 2048 {
		return errors.New("url de imagem muito longa")
	}
	return nil
}

func (servico *PerfilService) HistoricoPerfil(contexto context.Context, usuarioID string, limite int) (RespostaHistoricoPerfil, error) {
	if limite <= 0 || limite > 50 {
		limite = 20
	}
	return servico.repositorio.HistoricoPerfilUsuario(contexto, usuarioID, limite)
}

var (
	ErrTipoImagemInvalido = errors.New("tipo de imagem invalido")
	ErrCampoArquivoAusente = errors.New("campo file ou image obrigatorio")
)

// EnviarImagemPerfil processa, grava e persiste URL (sem recarregar organization_panel).
func (servico *PerfilService) EnviarImagemPerfil(contexto context.Context, usuarioID, perfilCodigoConta, tipo string, origem io.Reader) (RespostaUploadImagemPerfil, error) {
	tipo = strings.ToLower(strings.TrimSpace(tipo))
	var processado []byte
	var err error
	var subpasta string
	switch tipo {
	case "avatar":
		processado, err = media.ProcessarAvatar(origem)
		subpasta = "avatars"
	case "cover":
		processado, err = media.ProcessarCapa(origem)
		subpasta = "covers"
	default:
		return RespostaUploadImagemPerfil{}, ErrTipoImagemInvalido
	}
	if err != nil {
		return RespostaUploadImagemPerfil{}, err
	}
	caminhoRel, err := media.SalvarJPEG(media.ResolverDirUploads(), subpasta, media.NomeArquivoPerfil(usuarioID, tipo), processado)
	if err != nil {
		return RespostaUploadImagemPerfil{}, err
	}
	url := media.URLPublica(caminhoRel)
	if err := servico.repositorio.AtualizarURLImagemPerfil(contexto, usuarioID, perfilCodigoConta, tipo, url); err != nil {
		return RespostaUploadImagemPerfil{}, err
	}
	out := RespostaUploadImagemPerfil{}
	switch tipo {
	case "avatar":
		out.URLImagemAvatar = url
		out.URLAvatar = url
		out.URL = url
	case "cover":
		out.URLImagemCapa = url
		out.URLCapa = url
		out.URL = url
	}
	return out, nil
}
