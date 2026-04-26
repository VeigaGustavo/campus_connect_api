package comum

type PerfilPublicoAutor struct {
	Identificador string `json:"id"`
	Nome          string `json:"name"`
	URLAvatar     string `json:"avatar_image_url"`
	Perfil        string `json:"role"`
}
