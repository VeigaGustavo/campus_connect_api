package manipuladores

import (
	"net/http"
	"strings"

	"campus_connect_api/internal/modelos"
	"campus_connect_api/internal/respostas"
)

// VerificarSaude responde se o processo está no ar (liveness simples).
func VerificarSaude(resposta http.ResponseWriter, _ *http.Request) {
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "ok"})
}

// ListarFeedDescobrir GET /discover?filter=...
func ListarFeedDescobrir(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	filtro := requisicao.URL.Query().Get("filter")
	if filtro == "" {
		filtro = "all"
	}
	_ = filtro // TODO: aplicar filtro quando houver fonte de dados
	respostas.EscreverJSON(resposta, http.StatusOK, []modelos.ItemDescobrir{})
}

// ListarOportunidades GET /opportunities?q=...
func ListarOportunidades(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	_ = requisicao.URL.Query().Get("q")
	respostas.EscreverJSON(resposta, http.StatusOK, []modelos.Oportunidade{})
}

// ObterOportunidadePorIdentificador GET /opportunities/{id}
func ObterOportunidadePorIdentificador(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	identificador := strings.TrimPrefix(requisicao.URL.Path, "/opportunities/")
	if identificador == "" || identificador == "opportunities" {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "id_invalido", "identificador da oportunidade ausente")
		return
	}
	respostas.EscreverErro(resposta, http.StatusNotFound, "nao_encontrado", "oportunidade não encontrada")
}

// ListarEventos GET /events
func ListarEventos(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, []modelos.EventoCampus{})
}

// ObterEventoPorIdentificador GET /events/{id}
func ObterEventoPorIdentificador(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	identificador := strings.TrimPrefix(requisicao.URL.Path, "/events/")
	if identificador == "" || identificador == "events" {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "id_invalido", "identificador do evento ausente")
		return
	}
	respostas.EscreverErro(resposta, http.StatusNotFound, "nao_encontrado", "evento não encontrado")
}

// ListarGruposEstudo GET /groups
func ListarGruposEstudo(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, []modelos.GrupoEstudo{})
}

// ObterGrupoPorIdentificador GET /groups/{id} (stub até existir domínio).
func ObterGrupoPorIdentificador(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	identificador := strings.TrimPrefix(requisicao.URL.Path, "/groups/")
	identificador = strings.Trim(identificador, "/")
	if identificador == "" {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "id_invalido", "identificador do grupo ausente")
		return
	}
	respostas.EscreverErro(resposta, http.StatusNotFound, "nao_encontrado", "grupo não encontrado")
}

// SolicitarIngressoNoGrupo POST /groups/{id}/join
func SolicitarIngressoNoGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodPost {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use POST")
		return
	}
	resto := strings.TrimPrefix(requisicao.URL.Path, "/groups/")
	partes := strings.Split(strings.Trim(resto, "/"), "/")
	if len(partes) < 2 || partes[1] != "join" {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "caminho_invalido", "esperado /groups/{id}/join")
		return
	}
	_ = partes[0]
	respostas.EscreverErro(resposta, http.StatusNotImplemented, "nao_implementado", "ingresso exige autenticação e persistência")
}

// ObterPerfilDoUsuarioAutenticado GET /me ou /users/me
func ObterPerfilDoUsuarioAutenticado(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodGet {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use GET")
		return
	}
	respostas.EscreverErro(resposta, http.StatusUnauthorized, "nao_autorizado", "autenticação necessária")
}

// RegistrarCandidaturaEmOportunidade POST /opportunities/{id}/applications
func RegistrarCandidaturaEmOportunidade(resposta http.ResponseWriter, requisicao *http.Request) {
	if requisicao.Method != http.MethodPost {
		respostas.EscreverErro(resposta, http.StatusMethodNotAllowed, "metodo_nao_permitido", "use POST")
		return
	}
	const prefixoOportunidades = "/opportunities/"
	if !strings.HasPrefix(requisicao.URL.Path, prefixoOportunidades) || !strings.HasSuffix(requisicao.URL.Path, "/applications") {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "caminho_invalido", "esperado /opportunities/{id}/applications")
		return
	}
	identificadorOportunidade := strings.TrimSuffix(strings.TrimPrefix(requisicao.URL.Path, prefixoOportunidades), "/applications")
	identificadorOportunidade = strings.Trim(identificadorOportunidade, "/")
	if identificadorOportunidade == "" {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "id_invalido", "identificador da oportunidade ausente")
		return
	}
	respostas.EscreverErro(resposta, http.StatusNotImplemented, "nao_implementado", "fluxo de candidatura ainda não conectado")
}
