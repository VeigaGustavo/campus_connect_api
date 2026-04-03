package models

// DiscoverKind espelha DiscoverKind no app (feed Descobrir).
type DiscoverKind string

const (
	DiscoverKindInternship DiscoverKind = "internship"
	DiscoverKindEvent      DiscoverKind = "event"
	DiscoverKindStudyGroup DiscoverKind = "study_group"
	DiscoverKindProject    DiscoverKind = "project"
)

// WorkLocation — oportunidade (trabalho).
type WorkLocation string

const (
	WorkLocationRemote WorkLocation = "remote"
	WorkLocationHybrid WorkLocation = "hybrid"
	WorkLocationOnSite WorkLocation = "on_site"
)

// GroupLevel — grupo de estudo.
type GroupLevel string

const (
	GroupLevelBeginner     GroupLevel = "beginner"
	GroupLevelIntermediate GroupLevel = "intermediate"
	GroupLevelAdvanced     GroupLevel = "advanced"
)

// ProfileActivityKind — atividade recente no perfil.
type ProfileActivityKind string

const (
	ProfileActivityApplication     ProfileActivityKind = "application"
	ProfileActivityGroupJoined     ProfileActivityKind = "group_joined"
	ProfileActivityEventRegistered ProfileActivityKind = "event_registered"
)
