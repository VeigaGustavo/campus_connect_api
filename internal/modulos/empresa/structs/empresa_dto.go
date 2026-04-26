package model

type RequisicaoCriarOportunidade struct {
	Titulo            string   `json:"title"`
	NomeEmpresa       string   `json:"company_name"`
	DescricaoCurta    string   `json:"short_description"`
	DescricaoCompleta string   `json:"full_description"`
	PrazoCandidatura  string   `json:"apply_deadline"`
	ModalidadeLocal   string   `json:"work_location"`
	RotuloTipo        string   `json:"type_label"`
	Requisitos        []string `json:"requirements"`
	EscopoPublicacao  string   `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string   `json:"publish_group_id,omitempty"`
}

type CandidatoOportunidade struct {
	UsuarioID         string `json:"user_id"`
	Nome              string `json:"name"`
	Email             string `json:"email"`
	ResumoPerfil      string `json:"profile_summary"`
	URLCurriculo      string `json:"resume_url"`
	StatusCandidatura string `json:"application_status"`
}
