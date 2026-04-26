package model

type PerfilUsuario struct {
	Nome                        string                 `json:"name"`
	Iniciais                    string                 `json:"initials"`
	URLImagemCapa               string                 `json:"cover_image_url"`
	URLImagemAvatar             string                 `json:"avatar_image_url"`
	RotuloCertificadoDesempenho string                 `json:"performance_certificate_label"`
	CursoESemestre              string                 `json:"course_and_semester"`
	Email                       string                 `json:"email"`
	CidadeEstado                string                 `json:"city_state"`
	TotalCandidaturas           int                    `json:"applications_count"`
	TotalGrupos                 int                    `json:"groups_count"`
	TotalEventos                int                    `json:"events_count"`
	Interesses                  []InteressePerfil      `json:"interests"`
	AtividadesRecentes          []LinhaAtividadePerfil `json:"recent_activity"`
}
