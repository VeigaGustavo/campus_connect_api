package model

type DestaqueComunidadePerfil struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	Tipo          string `json:"kind"`
	Papel         string `json:"role"`
}

type ResumoPublicacaoOrganizacao struct {
	Identificador string `json:"id"`
	Titulo        string `json:"title"`
	Subtitulo     string `json:"subtitle"`
	IDReferencia  string `json:"reference_id"`
	CriadoEm      string `json:"created_at"`
}

type ResumoPostOrganizacao struct {
	Identificador string `json:"id"`
	Antevisao     string `json:"preview"`
	CriadoEm      string `json:"created_at"`
}

type PainelOrganizacaoPerfil struct {
	InstituicaoPai string `json:"parent_institution,omitempty"` // comunidade: instituicao do cadastro
	MapURL         string `json:"map_url,omitempty"`            // universidade: URL de mapa (cadastro / PUT)

	Vagas    []ResumoPublicacaoOrganizacao `json:"jobs,omitempty"`
	Eventos  []ResumoPublicacaoOrganizacao `json:"events,omitempty"`
	Grupos   []ResumoPublicacaoOrganizacao `json:"groups,omitempty"`
	Posts    []ResumoPostOrganizacao       `json:"posts,omitempty"`
	VagasTotal    int `json:"jobs_total"`
	EventosTotal  int `json:"events_total"`
	GruposTotal   int `json:"groups_total"`
	PostsTotal    int `json:"posts_total"`
}

type PerfilUsuario struct {
	ContextoPerfil     string                    `json:"profile_context"` // "user" | "organization"
	TipoPerfil         string                    `json:"profile_type"`    // estudante | comunidade | empresa | universidade
	PapelConta         string                    `json:"role"`              // padrao | comunidade | empresa | universidade | sistema_admin
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
	PainelOrganizacao  *PainelOrganizacaoPerfil   `json:"organization_panel,omitempty"`
}
