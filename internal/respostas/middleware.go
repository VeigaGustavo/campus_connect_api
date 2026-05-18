package respostas

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
)

type chaveContexto string

const chaveRequestID chaveContexto = "request_id"

func RequestIDFromContext(contexto context.Context) string {
	valorContexto := contexto.Value(chaveRequestID)
	s, _ := valorContexto.(string)
	return s
}

func EncadearRequestID(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(resposta http.ResponseWriter, requisicao *http.Request) {
		rid := strings.TrimSpace(requisicao.Header.Get("X-Request-Id"))
		if rid == "" {
			var bytesAleatorios [16]byte
			_, _ = rand.Read(bytesAleatorios[:])
			rid = hex.EncodeToString(bytesAleatorios[:])
		}
		resposta.Header().Set("X-Request-Id", rid)
		contexto := context.WithValue(requisicao.Context(), chaveRequestID, rid)
		proximo.ServeHTTP(resposta, requisicao.WithContext(contexto))
	})
}

func EncadearCORS(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(resposta http.ResponseWriter, requisicao *http.Request) {
		origem := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
		if origem == "" {
			origem = "*"
		}
		resposta.Header().Set("Access-Control-Allow-Origin", origem)
		resposta.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resposta.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if requisicao.Method == http.MethodOptions {
			resposta.WriteHeader(http.StatusNoContent)
			return
		}
		proximo.ServeHTTP(resposta, requisicao)
	})
}

func EncadearAceitarJSON(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(resposta http.ResponseWriter, requisicao *http.Request) {
		if requisicao.Method == http.MethodOptions {
			proximo.ServeHTTP(resposta, requisicao)
			return
		}
		aceitar := requisicao.Header.Get("Accept")
		if aceitar == "" || strings.Contains(aceitar, "*/*") || strings.Contains(aceitar, "application/json") {
			proximo.ServeHTTP(resposta, requisicao)
			return
		}
		EscreverErro(resposta, http.StatusNotAcceptable, "not_acceptable", "Accept must include application/json")
	})
}
