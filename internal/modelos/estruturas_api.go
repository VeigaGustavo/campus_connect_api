package modelos

// ItemDescobrir representa um cartão no feed Descobrir.
type ItemDescobrir struct {
	Identificador  string                 `json:"id"`
	Categoria      CategoriaItemDescobrir `json:"kind"`
	Titulo         string                 `json:"title"`
	Subtitulo      string                 `json:"subtitle,omitempty"`
	Resumo         string                 `json:"excerpt,omitempty"`
	MetaPrincipal  string                 `json:"meta_primary,omitempty"`
	MetaSecundaria string                 `json:"meta_secondary,omitempty"`
	IDReferencia   string                 `json:"reference_id"`
}

// Oportunidade lista ou detalhe de vaga/estágio.
type Oportunidade struct {
	Identificador     string                  `json:"id"`
	Titulo            string                  `json:"title"`
	NomeEmpresa       string                  `json:"company_name"`
	DescricaoCurta    string                  `json:"short_description"`
	DescricaoCompleta string                  `json:"full_description,omitempty"`
	PrazoCandidatura  string                  `json:"apply_deadline"` // ISO 8601
	ModalidadeLocal   ModalidadeLocalTrabalho `json:"work_location"`
	RotuloTipo        string                  `json:"type_label,omitempty"`
	Requisitos        []string                `json:"requirements,omitempty"`
}

// EventoCampus evento institucional.
type EventoCampus struct {
	Identificador string `json:"id"`
	Titulo        string `json:"title"`
	Descricao     string `json:"description,omitempty"`
	InicioEm      string `json:"start_at"` // ISO 8601
	Local         string `json:"location,omitempty"`
	Organizador   string `json:"organizer,omitempty"`
}

// GrupoEstudo grupo de estudos colaborativo.
type GrupoEstudo struct {
	Identificador string           `json:"id"`
	Titulo        string           `json:"title"`
	AreaEstudo    string           `json:"field_of_study"`
	Descricao     string           `json:"description,omitempty"`
	Nivel         NivelGrupoEstudo `json:"level"`
	TotalMembros  int              `json:"member_count"`
	RotuloHorario string           `json:"schedule_label,omitempty"`
}

// InteressePerfil interesse exibido no perfil.
type InteressePerfil struct {
	Rotulo string `json:"label"`
}

// LinhaAtividadePerfil uma entrada de atividade recente.
type LinhaAtividadePerfil struct {
	Tipo               TipoAtividadePerfil `json:"kind"`
	DestaqueTitulo     string              `json:"title_highlight"`
	Subtitulo          string              `json:"subtitle,omitempty"`
	TextoTempoRelativo string              `json:"time_ago_label,omitempty"`
	OcorridoEm         string              `json:"occurred_at,omitempty"` // ISO 8601
}

// PerfilUsuario resposta de GET /me ou /users/me.
type PerfilUsuario struct {
	Nome               string                 `json:"name"`
	Iniciais           string                 `json:"initials,omitempty"`
	CursoESemestre     string                 `json:"course_and_semester,omitempty"`
	Email              string                 `json:"email"`
	CidadeEstado       string                 `json:"city_state,omitempty"`
	TotalCandidaturas  int                    `json:"applications_count"`
	TotalGrupos        int                    `json:"groups_count"`
	TotalEventos       int                    `json:"events_count"`
	Interesses         []InteressePerfil      `json:"interests"`
	AtividadesRecentes []LinhaAtividadePerfil `json:"recent_activity"`
}

// ErroAPI corpo JSON padrão para erros HTTP.
type ErroAPI struct {
	Codigo   string `json:"code"`
	Mensagem string `json:"message"`
}
