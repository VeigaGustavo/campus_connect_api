package service

import "testing"

func TestParsePrazoCandidatura_FlutterSemFuso(t *testing.T) {
	casos := []string{
		"2026-08-31T23:59:59.000Z",
		"2026-08-31T23:59:59.000",
		"2026-06-30T00:00:00.000",
		"2026-08-31",
	}
	for _, s := range casos {
		if _, err := ParsePrazoCandidatura(s); err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
	}
}

func TestValidarRequisicaoOportunidade_OK(t *testing.T) {
	corpo := RequisicaoCriarOportunidade{
		Titulo:            "Dev",
		NomeEmpresa:       "Acme",
		DescricaoCurta:    "Resumo",
		DescricaoCompleta: "Longo",
		PrazoCandidatura:  "2026-08-31T23:59:59.000",
		ModalidadeLocal:   "hybrid",
		RotuloTipo:        "Estágio",
		Requisitos:        []string{"Flutter"},
	}
	if err := validarRequisicaoOportunidade(corpo); err != nil {
		t.Fatal(err)
	}
}

func TestValidarRequisicaoOportunidade_WorkLocationInvalido(t *testing.T) {
	corpo := RequisicaoCriarOportunidade{
		Titulo: "x", NomeEmpresa: "y", DescricaoCurta: "z", DescricaoCompleta: "w",
		PrazoCandidatura: "2026-08-31T23:59:59.000", ModalidadeLocal: "presencial",
		RotuloTipo: "Estágio", Requisitos: []string{},
	}
	if err := validarRequisicaoOportunidade(corpo); err == nil {
		t.Fatal("esperava erro work_location")
	}
}
