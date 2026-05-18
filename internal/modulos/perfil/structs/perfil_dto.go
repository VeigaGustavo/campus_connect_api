package model

type RequisicaoAtualizarPerfil struct {
	SobreMim         string   `json:"about_me"`
	Cargo            string   `json:"job_title"`
	Curso            string   `json:"course"`
	Semestre         string   `json:"semester"`
	Instituicao      string   `json:"institution_name"`
	MapURL           string   `json:"map_url,omitempty"`
	URLImagemAvatar  *string  `json:"avatar_image_url,omitempty"`
	URLImagemCapa    *string  `json:"cover_image_url,omitempty"`
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
	Itens []ItemHistoricoPerfil `json:"items"`
}

type RespostaUploadImagemPerfil struct {
	URLImagemAvatar string `json:"avatar_image_url,omitempty"`
	URLAvatar       string `json:"avatar_url,omitempty"`
	URL             string `json:"url,omitempty"`
	URLImagemCapa   string `json:"cover_image_url,omitempty"`
	URLCapa         string `json:"cover_url,omitempty"`
}
