package service

import (
	"errors"
	"fmt"
	"strings"
)

var ErrLeituraInvalida = errors.New("leitura invalida")

func normalizarTipoLeitura(kind string) (string, error) {
	k := strings.TrimSpace(strings.ToLower(kind))
	aliases := map[string]string{
		"revista":   string(LeituraRevista),
		"magazines": string(LeituraRevista),
		"artigo":    string(LeituraArtigo),
		"artigos":   string(LeituraArtigo),
		"articles":  string(LeituraArtigo),
		"noticia":   string(LeituraNoticiaCampus),
		"noticias":  string(LeituraNoticiaCampus),
		"news":      string(LeituraNoticiaCampus),
	}
	if canon, ok := aliases[k]; ok {
		k = canon
	}
	switch TipoItemLeitura(k) {
	case LeituraNoticiaCampus, LeituraRevista, LeituraArtigo:
		return k, nil
	default:
		return "", fmt.Errorf("%w: kind invalido (use campus_news, magazine ou article; recebido: %q)", ErrLeituraInvalida, kind)
	}
}

func validarRequisicaoLeituraSemanal(corpo *RequisicaoLeituraSemanal) error {
	tipo, err := normalizarTipoLeitura(corpo.Tipo)
	if err != nil {
		return err
	}
	corpo.Tipo = tipo
	if strings.TrimSpace(corpo.Titulo) == "" {
		return fmt.Errorf("%w: title obrigatorio", ErrLeituraInvalida)
	}
	if strings.TrimSpace(corpo.Fonte) == "" {
		return fmt.Errorf("%w: source obrigatorio", ErrLeituraInvalida)
	}
	if strings.TrimSpace(corpo.Resumo) == "" {
		return fmt.Errorf("%w: excerpt obrigatorio", ErrLeituraInvalida)
	}
	return nil
}
