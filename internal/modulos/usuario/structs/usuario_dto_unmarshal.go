package model

import (
	"encoding/json"
	"strconv"
	"strings"
)

func compactJSONKey(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if r != '_' && r != '-' && r != ' ' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// UnmarshalJSON aceita snake_case, camelCase, chaves com maiúsculas e cpf/cnpj como número JSON.
func (r *RequisicaoCadastroUsuario) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	norm := make(map[string]json.RawMessage, len(raw))
	for k, v := range raw {
		norm[compactJSONKey(k)] = v
	}
	get := func(keyCompact string) string {
		v, ok := norm[keyCompact]
		if !ok {
			return ""
		}
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			return strings.TrimSpace(s)
		}
		var n json.Number
		if err := json.Unmarshal(v, &n); err == nil {
			return strings.TrimSpace(n.String())
		}
		var f float64
		if err := json.Unmarshal(v, &f); err == nil {
			return strconv.FormatInt(int64(f), 10)
		}
		return ""
	}
	getInt := func(keyCompact string) int {
		v, ok := norm[keyCompact]
		if !ok {
			return 0
		}
		var i int
		if err := json.Unmarshal(v, &i); err == nil {
			return i
		}
		var f float64
		if err := json.Unmarshal(v, &f); err == nil {
			return int(f)
		}
		return 0
	}

	*r = RequisicaoCadastroUsuario{
		TipoPerfil:           get("profiletype"),
		NomeCompleto:         get("fullname"),
		Idade:                getInt("age"),
		DataNascimento:       get("birthdate"),
		CPF:                  get("cpf"),
		Instituicao:          get("institution"),
		Cidade:               get("city"),
		Estado:               get("state"),
		Email:                strings.ToLower(get("email")),
		Senha:                get("password"),
		TipoComunidade:       get("communitytype"),
		NomeComunidade:       get("communityname"),
		TituloGrupo:          get("grouptitle"),
		DescricaoGrupo:       get("groupdescription"),
		VisibilidadeGrupo:    get("groupvisibility"),
		NomeEmpresa:          get("companyname"),
		CNPJ:                 get("companycnpj"),
		DescricaoEmpresa:     get("companydescription"),
		NomeInstituicao:      get("institutionname"),
		SiglaInstituicao:     get("institutionacronym"),
		TipoInstituicao:      get("institutiontype"),
		DescricaoInstituicao: get("institutiondescription"),
	}
	return nil
}
