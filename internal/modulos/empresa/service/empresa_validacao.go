package service

import (
	"errors"
	"fmt"
	"strings"

	"campus_connect_api/internal/comum/horario"
	model "campus_connect_api/internal/modulos/empresa/structs"
)

var ErrOportunidadeInvalida = errors.New("oportunidade invalida")

const (
	maxRequisitosOportunidade   = 20
	maxCaracteresRequisito      = 50
	maxCaracteresDescricaoCurta = 500
)

func validarRequisicaoOportunidade(corpo RequisicaoCriarOportunidade) error {
	if strings.TrimSpace(corpo.Titulo) == "" {
		return fmt.Errorf("%w: title obrigatorio", ErrOportunidadeInvalida)
	}
	if strings.TrimSpace(corpo.NomeEmpresa) == "" {
		return fmt.Errorf("%w: company_name obrigatorio", ErrOportunidadeInvalida)
	}
	if strings.TrimSpace(corpo.DescricaoCurta) == "" {
		return fmt.Errorf("%w: short_description obrigatorio", ErrOportunidadeInvalida)
	}
	if len(corpo.DescricaoCurta) > maxCaracteresDescricaoCurta {
		return fmt.Errorf("%w: short_description max %d caracteres", ErrOportunidadeInvalida, maxCaracteresDescricaoCurta)
	}
	if strings.TrimSpace(corpo.DescricaoCompleta) == "" {
		return fmt.Errorf("%w: full_description obrigatorio", ErrOportunidadeInvalida)
	}
	if strings.TrimSpace(corpo.PrazoCandidatura) == "" {
		return fmt.Errorf("%w: apply_deadline obrigatorio", ErrOportunidadeInvalida)
	}
	if _, err := horario.ParseISO8601(corpo.PrazoCandidatura); err != nil {
		return fmt.Errorf("%w: apply_deadline invalido (use ISO-8601, ex. 2026-08-31T23:59:59Z ou 2026-08-31T23:59:59.000)", ErrOportunidadeInvalida)
	}
	wl := strings.TrimSpace(corpo.ModalidadeLocal)
	switch model.ModalidadeLocalTrabalho(wl) {
	case model.TrabalhoRemoto, model.TrabalhoHibrido, model.TrabalhoPresencial:
	default:
		return fmt.Errorf("%w: work_location deve ser remote, hybrid ou on_site", ErrOportunidadeInvalida)
	}
	if strings.TrimSpace(corpo.RotuloTipo) == "" {
		return fmt.Errorf("%w: type_label obrigatorio", ErrOportunidadeInvalida)
	}
	if corpo.Requisitos == nil {
		return fmt.Errorf("%w: requirements obrigatorio (pode ser [])", ErrOportunidadeInvalida)
	}
	if len(corpo.Requisitos) > maxRequisitosOportunidade {
		return fmt.Errorf("%w: requirements max %d itens", ErrOportunidadeInvalida, maxRequisitosOportunidade)
	}
	vistos := make(map[string]struct{}, len(corpo.Requisitos))
	for _, req := range corpo.Requisitos {
		req = strings.TrimSpace(req)
		if req == "" {
			continue
		}
		if len(req) > maxCaracteresRequisito {
			return fmt.Errorf("%w: cada requirement max %d caracteres", ErrOportunidadeInvalida, maxCaracteresRequisito)
		}
		chave := strings.ToLower(req)
		if _, ok := vistos[chave]; ok {
			continue
		}
		vistos[chave] = struct{}{}
	}
	return nil
}
