package service

import (
	"fmt"
	"net/url"
	"strings"
)

var tiposAnexoPermitidos = map[string]struct{}{
	"image": {}, "video": {}, "link": {}, "file": {},
}

var tiposConteudoPermitidos = map[string]struct{}{
	"article": {}, "campus_news": {}, "magazine": {}, "notice": {}, "project": {},
}

func validarCriarPost(corpo RequisicaoCriarPost) error {
	texto := strings.TrimSpace(corpo.Texto)
	if texto == "" && len(corpo.Anexos) == 0 {
		return fmt.Errorf("%w: texto ou anexo obrigatorio", ErrPostInvalido)
	}
	scope := strings.TrimSpace(corpo.EscopoPublicacao)
	if scope == "" {
		scope = "all"
	}
	if scope != "all" && scope != "group" {
		return fmt.Errorf("%w: publish_scope invalido", ErrPostInvalido)
	}
	if scope == "group" && strings.TrimSpace(corpo.IDGrupoPublicacao) == "" {
		return fmt.Errorf("%w: publish_group_id obrigatorio", ErrPostInvalido)
	}
	if k := strings.TrimSpace(corpo.TipoConteudo); k != "" {
		if _, ok := tiposConteudoPermitidos[k]; !ok {
			return fmt.Errorf("%w: content_kind invalido", ErrPostInvalido)
		}
	}
	for i, a := range corpo.Anexos {
		tipo := strings.ToLower(strings.TrimSpace(a.Tipo))
		if _, ok := tiposAnexoPermitidos[tipo]; !ok {
			return fmt.Errorf("%w: attachments[%d].type invalido", ErrPostInvalido, i)
		}
		u := strings.TrimSpace(a.URL)
		if u == "" {
			return fmt.Errorf("%w: attachments[%d].url obrigatorio", ErrPostInvalido, i)
		}
		if tipo == "link" || tipo == "file" {
			if err := validarURLExterna(u); err != nil {
				return fmt.Errorf("%w: attachments[%d].url invalida", ErrPostInvalido, i)
			}
		}
	}
	return nil
}

func validarURLExterna(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrPostInvalido
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrPostInvalido
	}
	return nil
}

func normalizarAnexos(anexos []AnexoPost) []AnexoPost {
	if len(anexos) == 0 {
		return []AnexoPost{}
	}
	out := make([]AnexoPost, 0, len(anexos))
	for _, a := range anexos {
		out = append(out, AnexoPost{
			Tipo: strings.ToLower(strings.TrimSpace(a.Tipo)),
			URL:  strings.TrimSpace(a.URL),
			Nome: strings.TrimSpace(a.Nome),
		})
	}
	return out
}
