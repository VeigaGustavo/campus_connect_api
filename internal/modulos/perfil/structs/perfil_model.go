package model

type DestaqueComunidadePerfil struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	Tipo          string `json:"kind"`
	Papel         string `json:"role"`
}

type PerfilUsuario struct {
	Identificador      string                    `json:"id"`
	Nome               string                    `json:"name"`
	URLImagemCapa      string                    `json:"cover_image_url"`
	URLImagemAvatar    string                    `json:"avatar_image_url"`
	Email              string                    `json:"email"`
	CidadeEstado       string                    `json:"city_state"`
	SobreMim           string                    `json:"about_me"`
	Cargo              string                    `json:"job_title"`
	Curso              string                    `json:"course"`
	Semestre           string                    `json:"semester"`
	Instituicao        string                    `json:"institution_name"`
	TotalCandidaturas  int                       `json:"applications_count"`
	TotalGrupos        int                       `json:"groups_count"`
	TotalEventos       int                       `json:"events_count"`
	Interesses         []string                  `json:"interests"`
	TopicosFavoritos   []string                  `json:"favorite_topics"`
	Especialidades     []string                  `json:"specialties"`
	DestaqueComunidade *DestaqueComunidadePerfil `json:"community_highlight"`
}
