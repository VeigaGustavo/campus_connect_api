package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	comum "campus_connect_api/internal/modulos/comum"
	grupoService "campus_connect_api/internal/modulos/grupo/service"
	auth "campus_connect_api/internal/modulos/seguranca/auth"
	"campus_connect_api/internal/respostas"
	"github.com/gin-gonic/gin"
)

type GrupoHTTPHandler struct {
	servicoGrupo *grupoService.GrupoService
	hubChat      *HubChatGrupo
}

func NovoGrupoHTTPHandler(servicoGrupo *grupoService.GrupoService) *GrupoHTTPHandler {
	return &GrupoHTTPHandler{
		servicoGrupo: servicoGrupo,
		hubChat:      NovoHubChatGrupo(),
	}
}

func (handler *GrupoHTTPHandler) RegistrarRotasGIN(grupo *gin.RouterGroup) {
	grupo.GET("/groups", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETGrupos)))
	grupo.POST("/groups", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTCriarGrupo)))
	grupo.PUT("/groups/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.PUTGrupo)))
	grupo.DELETE("/groups/:id", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.DELETEGrupo)))
	grupo.GET("/groups/:id/chat/messages", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETMensagensChatGrupo)))
	grupo.POST("/groups/:id/chat/messages", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTMensagemChatGrupo)))
	grupo.GET("/groups/:id/chat/ws", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.WSChatGrupo)))
	grupo.GET("/groups/:id/files", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETArquivosGrupo)))
	grupo.POST("/groups/:id/files", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTArquivoGrupo)))
	grupo.GET("/groups/:id/meetings", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETReunioesGrupo)))
	grupo.POST("/groups/:id/meetings", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTReuniaoGrupo)))
	grupo.GET("/groups/:id/events", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETEventosAssociadosGrupo)))
	grupo.POST("/groups/:id/events", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTAssociarEventoGrupo)))
	grupo.GET("/groups/:id/readings", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETLeiturasAssociadasGrupo)))
	grupo.POST("/groups/:id/readings", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTAssociarLeituraGrupo)))
	grupo.POST("/groups/:id/join", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTEntrarGrupo)))
	grupo.POST("/groups/:id/join-requests", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.POSTPedidoEntradaGrupo)))
	grupo.GET("/groups/:id/members", respostas.AdaptadorHTTP(auth.ExigirPerfis("padrao", "comunidade", "sistema_admin")(handler.GETMembrosGrupo)))
}

func (handler *GrupoHTTPHandler) GETGrupos(resposta http.ResponseWriter, requisicao *http.Request) {
	grupos, err := handler.servicoGrupo.ListarGrupos(requisicao.Context())
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if grupos == nil {
		grupos = []grupoService.GrupoEstudo{}
	}
	respostas.EscreverJSON(resposta, http.StatusOK, grupos)
}

func (handler *GrupoHTTPHandler) POSTCriarGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, ok := auth.SessaoDaRequisicao(requisicao)
	if !ok {
		respostas.EscreverErro(resposta, http.StatusUnauthorized, "unauthorized", "missing session")
		return
	}
	var corpo grupoService.RequisicaoCriarGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	grupoCriado, err := handler.servicoGrupo.CriarGrupo(requisicao.Context(), sessao.UsuarioID, corpo)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, grupoCriado)
}

func (handler *GrupoHTTPHandler) PUTGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoCriarGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	grupoAtualizado, err := handler.servicoGrupo.AtualizarGrupo(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil, corpo)
	if err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, grupoAtualizado)
}

func (handler *GrupoHTTPHandler) DELETEGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoGrupo.RemoverGrupo(requisicao.Context(), id, sessao.UsuarioID, sessao.Perfil); err != nil {
		respostas.EscreverErroPersistencia(resposta, err)
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "deleted"})
}

func (handler *GrupoHTTPHandler) GETMensagensChatGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	mensagens, err := handler.servicoGrupo.ListarMensagensGrupo(requisicao.Context(), id)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if mensagens == nil {
		mensagens = []grupoService.MensagemChatGrupo{}
	}
	respostas.EscreverJSON(resposta, http.StatusOK, mensagens)
}

func (handler *GrupoHTTPHandler) POSTMensagemChatGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoMensagemChatGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	mensagem, err := handler.servicoGrupo.AdicionarMensagemGrupo(requisicao.Context(), id, sessao.UsuarioID, corpo.Texto)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	handler.hubChat.EmitirMensagemChat(id, mensagem)
	respostas.EscreverJSON(resposta, http.StatusCreated, mensagem)
}

func (handler *GrupoHTTPHandler) WSChatGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	if err := handler.hubChat.Conectar(id, resposta, requisicao); err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "ws_error", err.Error())
		return
	}
}

func (handler *GrupoHTTPHandler) GETArquivosGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	respostas.EscreverJSON(resposta, http.StatusOK, handler.servicoGrupo.ListarArquivosGrupo(id))
}

func (handler *GrupoHTTPHandler) POSTArquivoGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoArquivoGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	arquivo := handler.servicoGrupo.AdicionarArquivoGrupo(id, sessao.UsuarioID, corpo.NomeArquivo, corpo.URLArquivo)
	respostas.EscreverJSON(resposta, http.StatusCreated, arquivo)
}

func (handler *GrupoHTTPHandler) GETReunioesGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	respostas.EscreverJSON(resposta, http.StatusOK, handler.servicoGrupo.ListarReunioesGrupo(id))
}

func (handler *GrupoHTTPHandler) POSTReuniaoGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoReuniaoGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	reuniao := handler.servicoGrupo.AdicionarReuniaoGrupo(id, corpo)
	respostas.EscreverJSON(resposta, http.StatusCreated, reuniao)
}

func (handler *GrupoHTTPHandler) GETEventosAssociadosGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	eventos := handler.servicoGrupo.ListarEventosAssociadosGrupo(id)
	if eventos == nil {
		eventos = []grupoService.AssociacaoGrupoEvento{}
	}
	respostas.EscreverJSON(resposta, http.StatusOK, eventos)
}

func (handler *GrupoHTTPHandler) POSTAssociarEventoGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoAssociarEventoGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	associacao := handler.servicoGrupo.AssociarEventoGrupo(id, corpo.EventoID)
	respostas.EscreverJSON(resposta, http.StatusCreated, associacao)
}

func (handler *GrupoHTTPHandler) GETLeiturasAssociadasGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	respostas.EscreverJSON(resposta, http.StatusOK, handler.servicoGrupo.ListarLeiturasAssociadasGrupo(id))
}

func (handler *GrupoHTTPHandler) POSTAssociarLeituraGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	var corpo grupoService.RequisicaoAssociarLeituraGrupo
	if err := json.NewDecoder(requisicao.Body).Decode(&corpo); err != nil {
		respostas.EscreverErro(resposta, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	associacao := handler.servicoGrupo.AssociarLeituraGrupo(id, corpo.LeituraID)
	respostas.EscreverJSON(resposta, http.StatusCreated, associacao)
}

func (handler *GrupoHTTPHandler) POSTEntrarGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoGrupo.EntrarGrupo(requisicao.Context(), id, sessao.UsuarioID); err != nil {
		switch {
		case errors.Is(err, comum.ErrNaoEncontrado):
			respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "group not found")
		case errors.Is(err, grupoService.ErrGrupoPrivado):
			respostas.EscreverErro(resposta, http.StatusForbidden, "private_group", "use join-requests for private groups")
		default:
			respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		}
		return
	}
	respostas.EscreverJSON(resposta, http.StatusOK, map[string]string{"status": "joined"})
}

func (handler *GrupoHTTPHandler) POSTPedidoEntradaGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	sessao, _ := auth.SessaoDaRequisicao(requisicao)
	id := requisicao.PathValue("id")
	if err := handler.servicoGrupo.PedirEntradaGrupo(requisicao.Context(), id, sessao.UsuarioID); err != nil {
		if errors.Is(err, comum.ErrNaoEncontrado) {
			respostas.EscreverErro(resposta, http.StatusNotFound, "not_found", "group not found")
			return
		}
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	respostas.EscreverJSON(resposta, http.StatusCreated, map[string]string{"status": "pending"})
}

func (handler *GrupoHTTPHandler) GETMembrosGrupo(resposta http.ResponseWriter, requisicao *http.Request) {
	id := requisicao.PathValue("id")
	membros, err := handler.servicoGrupo.ListarMembrosGrupo(requisicao.Context(), id)
	if err != nil {
		respostas.EscreverErro(resposta, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	if membros == nil {
		membros = []grupoService.MembroGrupo{}
	}
	respostas.EscreverJSON(resposta, http.StatusOK, membros)
}
