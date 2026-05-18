package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	feedService "campus_connect_api/internal/modulos/feed/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type FeedHTTPHandler struct {
	servicoFeed *feedService.FeedService
}

func NovoFeedHTTPHandler(servicoFeed *feedService.FeedService) *FeedHTTPHandler {
	return &FeedHTTPHandler{servicoFeed: servicoFeed}
}

func (handler *FeedHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/feed", respostas.AdaptadorHTTP(handler.GETFeed))
	grupo.POST("/feed/attachments", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "empresa", "universidade", "sistema_admin")(handler.POSTAnexoFeed)))
	grupo.GET("/feed/posts", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETListarPosts)))
	grupo.POST("/feed/posts", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "empresa", "universidade", "sistema_admin")(handler.POSTCriarPost)))
	grupo.GET("/feed/posts/:id", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETPost)))
	grupo.GET("/feed/posts/:id/comments", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.GETComentariosPost)))
	grupo.POST("/feed/posts/:id/comments", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.POSTComentarioPost)))
	grupo.PUT("/feed/posts/:id/reaction", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.PUTReacaoPost)))
	grupo.PUT("/feed/comments/:id/reaction", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.PUTReacaoComentario)))
	grupo.PUT("/feed/posts/:id/save", respostas.AdaptadorHTTP(auth.ObrigarAutenticacao(handler.PUTSalvarPost)))
}

func (handler *FeedHTTPHandler) GETFeed(resposta http.ResponseWriter, requisicao *http.Request) {
	filtro := requisicao.URL.Query().Get("filter")
	gruposCSV := strings.TrimSpace(requisicao.URL.Query().Get("group_ids"))
	var grupos []string
	if gruposCSV != "" {
		for _, g := range strings.Split(gruposCSV, ",") {
			gg := strings.TrimSpace(g)
			if gg != "" {
				grupos = append(grupos, gg)
			}
		}
	}
	out, err := handler.servicoFeed.Feed(requisicao.Context(), filtro, grupos)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *FeedHTTPHandler) GETListarPosts(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	pagina := 1
	if v := strings.TrimSpace(requisicao.URL.Query().Get("page")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pagina = n
		}
	}
	limite := 20
	if v := strings.TrimSpace(requisicao.URL.Query().Get("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limite = n
		}
	}
	incluirComentarios := strings.EqualFold(requisicao.URL.Query().Get("include_comments"), "true") ||
		requisicao.URL.Query().Get("include_comments") == "1"

	var grupos []string
	if gruposCSV := strings.TrimSpace(requisicao.URL.Query().Get("group_ids")); gruposCSV != "" {
		for _, g := range strings.Split(gruposCSV, ",") {
			gg := strings.TrimSpace(g)
			if gg != "" {
				grupos = append(grupos, gg)
			}
		}
	}

	out, err := handler.servicoFeed.ListarPosts(requisicao.Context(), sessao.UsuarioID, feedService.FiltroListarPosts{
		Pagina:             pagina,
		Limite:             limite,
		AutorID:            strings.TrimSpace(requisicao.URL.Query().Get("author_id")),
		GruposDoUsuario:    grupos,
		ApenasPostsDoGrupo: len(grupos) > 0,
		TipoConteudo:       strings.TrimSpace(requisicao.URL.Query().Get("content_kind")),
		IncluirComentarios: incluirComentarios,
	})
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, out)
}

func (handler *FeedHTTPHandler) POSTCriarPost(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo feedService.RequisicaoCriarPost
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	post, err := handler.servicoFeed.CriarPost(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		if errors.Is(err, feedService.ErrPostInvalido) {
			respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_post", err.Error())
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, post)
}

func (handler *FeedHTTPHandler) GETPost(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	postID := requisicao.PathValue("id")
	post, encontrado, err := handler.servicoFeed.ObterPost(requisicao.Context(), postID, sessao.UsuarioID)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if !encontrado {
		respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "resource not found")
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, post)
}

func (handler *FeedHTTPHandler) GETComentariosPost(resposta http.ResponseWriter, requisicao *http.Request) {
	postID := requisicao.PathValue("id")
	out, err := handler.servicoFeed.ListarComentariosPost(requisicao.Context(), postID)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]any{"items": out})
}

func (handler *FeedHTTPHandler) POSTComentarioPost(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	postID := requisicao.PathValue("id")
	var corpo feedService.RequisicaoCriarComentario
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	out, err := handler.servicoFeed.CriarComentarioPost(requisicao.Context(), postID, sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, out)
}

func (handler *FeedHTTPHandler) PUTReacaoPost(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	postID := requisicao.PathValue("id")
	var corpo feedService.RequisicaoReacao
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if err := handler.servicoFeed.ReagirPost(requisicao.Context(), postID, sessao.UsuarioID, corpo.Reacao); err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "ok"})
}

func (handler *FeedHTTPHandler) PUTReacaoComentario(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	comentarioID := requisicao.PathValue("id")
	var corpo feedService.RequisicaoReacao
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if err := handler.servicoFeed.ReagirComentario(requisicao.Context(), comentarioID, sessao.UsuarioID, corpo.Reacao); err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "ok"})
}

func (handler *FeedHTTPHandler) PUTSalvarPost(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	postID := requisicao.PathValue("id")
	var corpo feedService.RequisicaoSalvarPost
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if err := handler.servicoFeed.SalvarPost(requisicao.Context(), postID, sessao.UsuarioID, corpo.Salvo); err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "ok"})
}
