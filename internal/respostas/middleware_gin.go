package respostas

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GinRequestID() gin.HandlerFunc {
	return func(contexto *gin.Context) {
		rid := strings.TrimSpace(contexto.GetHeader("X-Request-Id"))
		if rid == "" {
			var bytesAleatorios [16]byte
			_, _ = rand.Read(bytesAleatorios[:])
			rid = hex.EncodeToString(bytesAleatorios[:])
		}
		contexto.Writer.Header().Set("X-Request-Id", rid)

		contextoRequisicao := context.WithValue(contexto.Request.Context(), chaveRequestID, rid)
		contexto.Request = contexto.Request.WithContext(contextoRequisicao)
		contexto.Next()
		if codigo := contexto.Writer.Status(); codigo >= http.StatusInternalServerError {
			slog.Default().Error("http_response_error",
				"method", contexto.Request.Method,
				"path", contexto.Request.URL.Path,
				"status", codigo,
				"request_id", rid,
			)
		}
	}
}

func GinCORS() gin.HandlerFunc {
	return func(contexto *gin.Context) {
		origem := strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
		if origem == "" {
			origem = "*"
		}
		contexto.Writer.Header().Set("Access-Control-Allow-Origin", origem)
		contexto.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		contexto.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

		if contexto.Request.Method == http.MethodOptions {
			contexto.Status(http.StatusNoContent)
			contexto.Abort()
			return
		}
		contexto.Next()
	}
}

func GinAceitarJSON() gin.HandlerFunc {
	return func(contexto *gin.Context) {
		if contexto.Request.Method == http.MethodOptions {
			contexto.Next()
			return
		}
		// Upload de imagem: multipart, nao exige Accept JSON.
		if contexto.Request.Method == http.MethodPost {
			p := contexto.Request.URL.Path
			if strings.HasSuffix(p, "/profile/avatar") || strings.HasSuffix(p, "/profile/cover") || strings.HasSuffix(p, "/feed/attachments") {
				contexto.Next()
				return
			}
		}
		aceitar := contexto.GetHeader("Accept")
		if aceitar == "" || strings.Contains(aceitar, "*/*") || strings.Contains(aceitar, "application/json") {
			contexto.Next()
			return
		}
		EscreverErro(contexto.Writer, http.StatusNotAcceptable, "not_acceptable", "Accept must include application/json")
		contexto.Abort()
	}
}
