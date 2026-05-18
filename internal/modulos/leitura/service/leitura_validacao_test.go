package service

import "testing"

func TestNormalizarTipoLeitura_Revista(t *testing.T) {
	got, err := normalizarTipoLeitura("revista")
	if err != nil {
		t.Fatal(err)
	}
	if got != string(LeituraRevista) {
		t.Fatalf("got %q want magazine", got)
	}
}

func TestNormalizarTipoLeitura_Magazine(t *testing.T) {
	got, err := normalizarTipoLeitura("magazine")
	if err != nil || got != "magazine" {
		t.Fatalf("magazine: got %q err %v", got, err)
	}
}

func TestValidarRequisicaoLeituraSemanal_KindInvalido(t *testing.T) {
	corpo := RequisicaoLeituraSemanal{Tipo: "revistas", Titulo: "t", Fonte: "f", Resumo: "e"}
	if err := validarRequisicaoLeituraSemanal(&corpo); err == nil {
		t.Fatal("esperava erro")
	}
}
