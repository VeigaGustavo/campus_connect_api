package respostas

import (
	"encoding/json"
	"net/http"

	"campus_connect_api/internal/modelos"
	"campus_connect_api/internal/versao"
)

func EscreverJSON(resposta http.ResponseWriter, status int, corpo any) {
	resposta.Header().Set("Content-Type", "application/json; charset=utf-8")
	resposta.WriteHeader(status)
	_ = json.NewEncoder(resposta).Encode(corpo)
}

func EscreverErro(resposta http.ResponseWriter, status int, codigo, mensagem string) {
	EscreverJSON(resposta, status, modelos.ErroAPI{Codigo: codigo, Mensagem: mensagem, Revisao: versao.Revisao})
}
