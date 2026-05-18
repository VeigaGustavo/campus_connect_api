package model

type RequisicaoCriarGrupo struct {
	Titulo            string `json:"title"`
	AreaEstudo        string `json:"field_of_study"`
	Descricao         string `json:"description"`
	Nivel             string `json:"level"`
	RotuloHorario     string `json:"schedule_label"`
	Visibilidade      string `json:"visibility,omitempty"`
	EscopoPublicacao  string `json:"publish_scope,omitempty"`
	IDGrupoPublicacao string `json:"publish_group_id,omitempty"`
}

type RequisicaoMensagemChatGrupo struct {
	Texto string `json:"text"`
}

type RequisicaoArquivoGrupo struct {
	NomeArquivo string `json:"file_name"`
	URLArquivo  string `json:"file_url"`
}

type RequisicaoReuniaoGrupo struct {
	Tema          string `json:"topic"`
	InicioEm      string `json:"start_at"`
	Local         string `json:"location"`
	Participantes int    `json:"participants_count"`
}

type RequisicaoAssociarEventoGrupo struct {
	EventoID string `json:"event_id"`
}

type RequisicaoAssociarLeituraGrupo struct {
	LeituraID string `json:"reading_id"`
}
