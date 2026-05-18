package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	grupoService "campus_connect_api/internal/modulos/grupo/service"
	"github.com/gorilla/websocket"
)

type HubChatGrupo struct {
	mutex    sync.RWMutex
	salas    map[string]map[*websocket.Conn]struct{}
	upgrader websocket.Upgrader
}

func NovoHubChatGrupo() *HubChatGrupo {
	return &HubChatGrupo{
		salas: map[string]map[*websocket.Conn]struct{}{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(requisicao *http.Request) bool { return true },
		},
	}
}

func (hub *HubChatGrupo) Conectar(grupoID string, resposta http.ResponseWriter, requisicao *http.Request) error {
	conexao, err := hub.upgrader.Upgrade(resposta, requisicao, nil)
	if err != nil {
		return err
	}

	hub.mutex.Lock()
	if _, existe := hub.salas[grupoID]; !existe {
		hub.salas[grupoID] = map[*websocket.Conn]struct{}{}
	}
	hub.salas[grupoID][conexao] = struct{}{}
	hub.mutex.Unlock()

	go func() {
		defer func() {
			hub.mutex.Lock()
			delete(hub.salas[grupoID], conexao)
			hub.mutex.Unlock()
			_ = conexao.Close()
		}()
		for {
			tipoMensagem, conteudo, errLeitura := conexao.ReadMessage()
			if errLeitura != nil {
				return
			}
			hub.broadcast(grupoID, tipoMensagem, conteudo)
		}
	}()

	return nil
}

func (hub *HubChatGrupo) broadcast(grupoID string, tipoMensagem int, conteudo []byte) {
	hub.mutex.RLock()
	conexoes := hub.salas[grupoID]
	hub.mutex.RUnlock()
	for conexao := range conexoes {
		_ = conexao.WriteMessage(tipoMensagem, conteudo)
	}
}

func (hub *HubChatGrupo) EmitirMensagemChat(grupoID string, mensagem grupoService.MensagemChatGrupo) {
	payload, err := json.Marshal(map[string]any{
		"event":   "chat_message",
		"message": mensagem,
	})
	if err != nil {
		return
	}
	hub.broadcast(grupoID, websocket.TextMessage, payload)
}
