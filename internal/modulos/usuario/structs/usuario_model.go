package model

type Usuario struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	Email         string `json:"email"`
	Perfil        string `json:"role"`
}
