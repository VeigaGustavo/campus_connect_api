package model

type RequisicaoCriarUsuario struct {
	Nome   string `json:"name"`
	Email  string `json:"email"`
	Senha  string `json:"password"`
	Perfil string `json:"role"`
}

type RequisicaoCadastroUsuario struct {
	TipoPerfil           string `json:"profile_type"`
	NomeCompleto         string `json:"full_name"`
	Idade                int    `json:"age"`
	DataNascimento       string `json:"birth_date,omitempty"`
	CPF                  string `json:"cpf"`
	Instituicao          string `json:"institution"`
	Cidade               string `json:"city"`
	Estado               string `json:"state"`
	Email                string `json:"email"`
	Senha                string `json:"password"`
	TipoComunidade       string `json:"community_type,omitempty"`
	NomeComunidade       string `json:"community_name,omitempty"`
	TituloGrupo          string `json:"group_title,omitempty"`
	DescricaoGrupo       string `json:"group_description,omitempty"`
	VisibilidadeGrupo    string `json:"group_visibility,omitempty"`
	NomeEmpresa          string `json:"company_name,omitempty"`
	CNPJ                 string `json:"company_cnpj,omitempty"`
	DescricaoEmpresa     string `json:"company_description,omitempty"`
	NomeInstituicao      string `json:"institution_name,omitempty"`
	SiglaInstituicao     string `json:"institution_acronym,omitempty"`
	TipoInstituicao      string `json:"institution_type,omitempty"`
	DescricaoInstituicao string `json:"institution_description,omitempty"`
}
