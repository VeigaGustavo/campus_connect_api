package model

type InteressePerfil struct {
	Rotulo string `json:"label"`
}

type RequisicaoAtualizarPerfil struct {
	SobreMim         string   `json:"about_me"`
	Cargo            string   `json:"job_title"`
	Curso            string   `json:"course"`
	Semestre         string   `json:"semester"`
	Instituicao      string   `json:"institution_name"`
	Interesses       []string `json:"interests"`
	TopicosFavoritos []string `json:"favorite_topics"`
	Especialidades   []string `json:"specialties"`
}

type ItemHistoricoPerfil struct {
	Identificador string `json:"id"`
	Tipo          string `json:"kind"`
	Titulo        string `json:"title"`
	Subtitulo     string `json:"subtitle"`
	IDReferencia  string `json:"reference_id"`
	CriadoEm      string `json:"created_at"`
}

type RespostaHistoricoPerfil struct {
	Itens         []ItemHistoricoPerfil `json:"items"`
	ProximoCursor string                `json:"next_cursor,omitempty"`
}

type LinhaAtividadePerfil struct {
	Tipo               TipoAtividadePerfil `json:"kind"`
	DestaqueTitulo     string              `json:"title_highlight"`
	Subtitulo          string              `json:"subtitle"`
	TextoTempoRelativo string              `json:"time_ago_label"`
}
