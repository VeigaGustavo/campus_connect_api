package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"campus_connect_api/internal/manipuladores"
)

func main() {
	roteador := http.NewServeMux()
	roteador.HandleFunc("/health", manipuladores.VerificarSaude)

	roteador.HandleFunc("/discover", manipuladores.ListarFeedDescobrir)

	roteador.HandleFunc("/opportunities/", func(resposta http.ResponseWriter, requisicao *http.Request) {
		if requisicao.URL.Path == "/opportunities/" || requisicao.URL.Path == "/opportunities" {
			if requisicao.Method == http.MethodGet {
				manipuladores.ListarOportunidades(resposta, requisicao)
				return
			}
		}
		if strings.HasPrefix(requisicao.URL.Path, "/opportunities/") && strings.HasSuffix(requisicao.URL.Path, "/applications") {
			manipuladores.RegistrarCandidaturaEmOportunidade(resposta, requisicao)
			return
		}
		manipuladores.ObterOportunidadePorIdentificador(resposta, requisicao)
	})
	roteador.HandleFunc("/opportunities", manipuladores.ListarOportunidades)

	roteador.HandleFunc("/events/", func(resposta http.ResponseWriter, requisicao *http.Request) {
		if requisicao.URL.Path == "/events/" || requisicao.URL.Path == "/events" {
			manipuladores.ListarEventos(resposta, requisicao)
			return
		}
		manipuladores.ObterEventoPorIdentificador(resposta, requisicao)
	})
	roteador.HandleFunc("/events", manipuladores.ListarEventos)

	roteador.HandleFunc("/groups/", func(resposta http.ResponseWriter, requisicao *http.Request) {
		caminho := strings.TrimSuffix(requisicao.URL.Path, "/")
		if strings.HasSuffix(caminho, "/join") {
			manipuladores.SolicitarIngressoNoGrupo(resposta, requisicao)
			return
		}
		manipuladores.ObterGrupoPorIdentificador(resposta, requisicao)
	})
	roteador.HandleFunc("/groups", manipuladores.ListarGruposEstudo)

	roteador.HandleFunc("/me", manipuladores.ObterPerfilDoUsuarioAutenticado)
	roteador.HandleFunc("/users/me", manipuladores.ObterPerfilDoUsuarioAutenticado)

	enderecoEscuta := resolverEnderecoEscuta()
	log.Printf("campus_connect_api escutando em %s", enderecoEscuta)
	if err := http.ListenAndServe(enderecoEscuta, roteador); err != nil {
		log.Fatal(err)
	}
}

// resolverEnderecoEscuta: LISTEN_ADDR tem prioridade; senão PORT; padrão ":8080".
func resolverEnderecoEscuta() string {
	if definido := os.Getenv("LISTEN_ADDR"); definido != "" {
		return definido
	}
	if porta := os.Getenv("PORT"); porta != "" {
		return ":" + porta
	}
	return ":8080"
}
