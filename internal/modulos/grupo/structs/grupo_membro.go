package model

type MembroGrupo struct {
	UsuarioID  string `json:"user_id"`
	Nome       string `json:"name"`
	URLAvatar  string `json:"avatar_url"`
	Papel      string `json:"role"` // owner | admin | member
	EntrouEm   string `json:"joined_at,omitempty"`
}
