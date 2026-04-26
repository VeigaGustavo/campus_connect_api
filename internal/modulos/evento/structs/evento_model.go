package model

type EventoCampus struct {
	Identificador string `json:"id"`
	Titulo        string `json:"title"`
	Descricao     string `json:"description"`
	InicioEm      string `json:"start_at"`
	Local         string `json:"location"`
	Organizador   string `json:"organizer"`
}
