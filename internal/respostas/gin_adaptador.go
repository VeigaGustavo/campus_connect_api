package respostas

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdaptadorHTTP(handler http.HandlerFunc) gin.HandlerFunc {
	return func(contexto *gin.Context) {
		for _, parametro := range contexto.Params {
			contexto.Request.SetPathValue(parametro.Key, parametro.Value)
		}
		handler(contexto.Writer, contexto.Request)
	}
}
