package service

import (
	"encoding/json"
	"strings"
	"testing"

	model "campus_connect_api/internal/modulos/usuario/structs"
)

func TestDecodeRegisterEmpresa_JSON(t *testing.T) {
	raw := `{
  "profile_type": "empresa",
  "full_name": "Gustavo",
  "birth_date": "2006-05-15",
  "cpf": "06927249150",
  "city": "Palmas",
  "company_name": "Veiga.dev",
  "email": "gustavoavdcarmo@outlook.com",
  "password": "Veiga.2004",
  "state": "TO"
}`
	var corpo RequisicaoCadastroUsuario
	if err := json.Unmarshal([]byte(raw), &corpo); err != nil {
		t.Fatal(err)
	}
	if corpo.NomeEmpresa != "Veiga.dev" {
		t.Fatalf("company_name: got %q", corpo.NomeEmpresa)
	}
	if corpo.TipoPerfil != "empresa" {
		t.Fatalf("profile_type: got %q", corpo.TipoPerfil)
	}

	// Simula o mesmo fluxo do serviço (sem DB)
	corpo.TipoPerfil = strings.ToLower(strings.TrimSpace(corpo.TipoPerfil))
	corpo.DataNascimento = normalizarDataNascimento(corpo.DataNascimento)
	if corpo.Idade <= 0 {
		if idade, ok := idadeAPartirDeDataNascimento(strings.TrimSpace(corpo.DataNascimento)); ok {
			corpo.Idade = idade
		}
	}
	if corpo.Idade <= 0 {
		t.Fatalf("idade=%d birth=%q", corpo.Idade, corpo.DataNascimento)
	}
	if strings.TrimSpace(corpo.NomeEmpresa) == "" {
		t.Fatal("nome empresa vazio")
	}
}

func TestRequisicaoCadastroUsuario_UnmarshalJSON_CompanyNameCamel(t *testing.T) {
	raw := `{"profileType":"empresa","fullName":"X","companyName":"Acme","cpf":"1","city":"c","state":"s","email":"a@b.c","password":"p","birthDate":"2006-05-15"}`
	var corpo model.RequisicaoCadastroUsuario
	if err := json.Unmarshal([]byte(raw), &corpo); err != nil {
		t.Fatal(err)
	}
	if corpo.NomeEmpresa != "Acme" || corpo.TipoPerfil != "empresa" {
		t.Fatalf("%+v", corpo)
	}
}
