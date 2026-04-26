package model

type GrupoEstudo struct {
	Identificador string           `json:"id"`
	Titulo        string           `json:"title"`
	AreaEstudo    string           `json:"field_of_study"`
	Descricao     string           `json:"description"`
	Nivel         NivelGrupoEstudo `json:"level"`
	TotalMembros  int              `json:"member_count"`
	RotuloHorario string           `json:"schedule_label"`
}

type MensagemChatGrupo struct {
	ID       string `json:"id"`
	GrupoID  string `json:"group_id"`
	AutorID  string `json:"author_id"`
	Texto    string `json:"text"`
	CriadoEm string `json:"created_at"`
}

type ArquivoGrupo struct {
	ID          string `json:"id"`
	GrupoID     string `json:"group_id"`
	NomeArquivo string `json:"file_name"`
	URLArquivo  string `json:"file_url"`
	AutorID     string `json:"author_id"`
	CriadoEm    string `json:"created_at"`
}

type ReuniaoGrupo struct {
	ID            string `json:"id"`
	GrupoID       string `json:"group_id"`
	Tema          string `json:"topic"`
	InicioEm      string `json:"start_at"`
	Local         string `json:"location"`
	Participantes int    `json:"participants_count"`
}

type AssociacaoGrupoEvento struct {
	ID       string `json:"id"`
	GrupoID  string `json:"group_id"`
	EventoID string `json:"event_id"`
}

type AssociacaoGrupoLeitura struct {
	ID        string `json:"id"`
	GrupoID   string `json:"group_id"`
	LeituraID string `json:"reading_id"`
}
