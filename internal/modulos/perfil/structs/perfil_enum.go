package model

type TipoAtividadePerfil string

const (
	AtividadeCandidatura       TipoAtividadePerfil = "application"
	AtividadeEntradaEmGrupo    TipoAtividadePerfil = "group_joined"
	AtividadeInscricaoEmEvento TipoAtividadePerfil = "event_registered"
)
