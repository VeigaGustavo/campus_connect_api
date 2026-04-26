package model

type InteressePerfil struct {
	Rotulo string `json:"label"`
}

type LinhaAtividadePerfil struct {
	Tipo               TipoAtividadePerfil `json:"kind"`
	DestaqueTitulo     string              `json:"title_highlight"`
	Subtitulo          string              `json:"subtitle"`
	TextoTempoRelativo string              `json:"time_ago_label"`
}
