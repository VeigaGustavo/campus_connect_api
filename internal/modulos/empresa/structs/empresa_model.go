package model

type Oportunidade struct {
	Identificador     string                  `json:"id"`
	Titulo            string                  `json:"title"`
	NomeEmpresa       string                  `json:"company_name"`
	DescricaoCurta    string                  `json:"short_description"`
	DescricaoCompleta string                  `json:"full_description"`
	PrazoCandidatura  string                  `json:"apply_deadline"`
	ModalidadeLocal   ModalidadeLocalTrabalho `json:"work_location"`
	RotuloTipo        string                  `json:"type_label"`
	Requisitos        []string                `json:"requirements"`
}
