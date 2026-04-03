package modelos

// CategoriaItemDescobrir espelha o tipo de item no feed Descobrir do app (valores JSON estáveis).
type CategoriaItemDescobrir string

const (
	CategoriaEstagio     CategoriaItemDescobrir = "internship"
	CategoriaEvento      CategoriaItemDescobrir = "event"
	CategoriaGrupoEstudo CategoriaItemDescobrir = "study_group"
	CategoriaProjeto     CategoriaItemDescobrir = "project"
)

// ModalidadeLocalTrabalho descreve onde a vaga pode ser exercida (valores JSON estáveis).
type ModalidadeLocalTrabalho string

const (
	TrabalhoRemoto     ModalidadeLocalTrabalho = "remote"
	TrabalhoHibrido    ModalidadeLocalTrabalho = "hybrid"
	TrabalhoPresencial ModalidadeLocalTrabalho = "on_site"
)

// NivelGrupoEstudo dificuldade ou nível do grupo (valores JSON estáveis).
type NivelGrupoEstudo string

const (
	NivelIniciante     NivelGrupoEstudo = "beginner"
	NivelIntermediario NivelGrupoEstudo = "intermediate"
	NivelAvancado      NivelGrupoEstudo = "advanced"
)

// TipoAtividadePerfil tipo de linha no histórico do perfil (valores JSON estáveis).
type TipoAtividadePerfil string

const (
	AtividadeCandidatura     TipoAtividadePerfil = "application"
	AtividadeEntradaEmGrupo  TipoAtividadePerfil = "group_joined"
	AtividadeInscricaoEvento TipoAtividadePerfil = "event_registered"
)
