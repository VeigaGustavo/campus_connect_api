package modelos

// ErroAPI corpo JSON para erros HTTP.
type ErroAPI struct {
	Codigo   string `json:"code"`
	Mensagem string `json:"message"`
}

type RequisicaoLogin struct {
	Email string `json:"email"`
	Senha string `json:"password"`
}

type RespostaLogin struct {
	Token     string `json:"access_token"`
	TipoToken string `json:"token_type"`
	Login     string `json:"login"`
	ExpiraEm  int64  `json:"expires_in"`
	Perfil    string `json:"role"`
	UsuarioID string `json:"user_id"`
}

type UsuarioSessao struct {
	ID     string `json:"id"`
	Nome   string `json:"name,omitempty"`
	Email  string `json:"email"`
	Perfil string `json:"role"`
}
