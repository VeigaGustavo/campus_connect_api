package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"campus_connect_api/internal/infra/database"
	grupoService "campus_connect_api/internal/modulos/grupo/service"
)

type grupoRepositoryPostgres struct {
	store *database.Postgres

	mutex              sync.RWMutex
	chatGrupo          map[string][]grupoService.MensagemChatGrupo
	arquivosGrupo      map[string][]grupoService.ArquivoGrupo
	reunioesGrupo      map[string][]grupoService.ReuniaoGrupo
	eventosAssociados  map[string][]grupoService.AssociacaoGrupoEvento
	leiturasAssociadas map[string][]grupoService.AssociacaoGrupoLeitura
}

func NovoGrupoRepository(store *database.Postgres) grupoService.GrupoRepository {
	return &grupoRepositoryPostgres{
		store:              store,
		chatGrupo:          map[string][]grupoService.MensagemChatGrupo{},
		arquivosGrupo:      map[string][]grupoService.ArquivoGrupo{},
		reunioesGrupo:      map[string][]grupoService.ReuniaoGrupo{},
		eventosAssociados:  map[string][]grupoService.AssociacaoGrupoEvento{},
		leiturasAssociadas: map[string][]grupoService.AssociacaoGrupoLeitura{},
	}
}

func novoIdentificador(prefixo string) string {
	var bytesAleatorios [8]byte
	_, _ = rand.Read(bytesAleatorios[:])
	return prefixo + hex.EncodeToString(bytesAleatorios[:])
}

func (repositorio *grupoRepositoryPostgres) ListarGrupos(contexto context.Context) ([]grupoService.GrupoEstudo, error) {
	return repositorio.store.ListarGrupos(contexto)
}

func (repositorio *grupoRepositoryPostgres) InserirGrupo(contexto context.Context, criadoPor string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return repositorio.store.InserirGrupo(contexto, criadoPor, corpo)
}

func (repositorio *grupoRepositoryPostgres) AtualizarGrupo(contexto context.Context, id, usuarioID string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return repositorio.store.AtualizarGrupo(contexto, id, usuarioID, corpo)
}

func (repositorio *grupoRepositoryPostgres) AtualizarGrupoComoAdmin(contexto context.Context, id string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return repositorio.store.AtualizarGrupoComoAdmin(contexto, id, corpo)
}

func (repositorio *grupoRepositoryPostgres) RemoverGrupo(contexto context.Context, id, usuarioID string) error {
	return repositorio.store.RemoverGrupo(contexto, id, usuarioID)
}

func (repositorio *grupoRepositoryPostgres) RemoverGrupoComoAdmin(contexto context.Context, id string) error {
	return repositorio.store.RemoverGrupoComoAdmin(contexto, id)
}

func (repositorio *grupoRepositoryPostgres) ListarMensagensGrupo(grupoID string) []grupoService.MensagemChatGrupo {
	repositorio.mutex.RLock()
	defer repositorio.mutex.RUnlock()
	return append([]grupoService.MensagemChatGrupo(nil), repositorio.chatGrupo[grupoID]...)
}

func (repositorio *grupoRepositoryPostgres) AdicionarMensagemGrupo(grupoID, autorID, texto string) grupoService.MensagemChatGrupo {
	mensagem := grupoService.MensagemChatGrupo{
		ID:       novoIdentificador("msg-"),
		GrupoID:  grupoID,
		AutorID:  autorID,
		Texto:    texto,
		CriadoEm: time.Now().UTC().Format(time.RFC3339),
	}
	repositorio.mutex.Lock()
	repositorio.chatGrupo[grupoID] = append(repositorio.chatGrupo[grupoID], mensagem)
	repositorio.mutex.Unlock()
	return mensagem
}

func (repositorio *grupoRepositoryPostgres) ListarArquivosGrupo(grupoID string) []grupoService.ArquivoGrupo {
	repositorio.mutex.RLock()
	defer repositorio.mutex.RUnlock()
	return append([]grupoService.ArquivoGrupo(nil), repositorio.arquivosGrupo[grupoID]...)
}

func (repositorio *grupoRepositoryPostgres) AdicionarArquivoGrupo(grupoID, autorID, nome, url string) grupoService.ArquivoGrupo {
	arquivo := grupoService.ArquivoGrupo{
		ID:          novoIdentificador("file-"),
		GrupoID:     grupoID,
		NomeArquivo: nome,
		URLArquivo:  url,
		AutorID:     autorID,
		CriadoEm:    time.Now().UTC().Format(time.RFC3339),
	}
	repositorio.mutex.Lock()
	repositorio.arquivosGrupo[grupoID] = append(repositorio.arquivosGrupo[grupoID], arquivo)
	repositorio.mutex.Unlock()
	return arquivo
}

func (repositorio *grupoRepositoryPostgres) ListarReunioesGrupo(grupoID string) []grupoService.ReuniaoGrupo {
	repositorio.mutex.RLock()
	defer repositorio.mutex.RUnlock()
	return append([]grupoService.ReuniaoGrupo(nil), repositorio.reunioesGrupo[grupoID]...)
}

func (repositorio *grupoRepositoryPostgres) AdicionarReuniaoGrupo(grupoID string, corpo grupoService.RequisicaoReuniaoGrupo) grupoService.ReuniaoGrupo {
	reuniao := grupoService.ReuniaoGrupo{
		ID:            novoIdentificador("meet-"),
		GrupoID:       grupoID,
		Tema:          corpo.Tema,
		InicioEm:      corpo.InicioEm,
		Local:         corpo.Local,
		Participantes: corpo.Participantes,
	}
	repositorio.mutex.Lock()
	repositorio.reunioesGrupo[grupoID] = append(repositorio.reunioesGrupo[grupoID], reuniao)
	repositorio.mutex.Unlock()
	return reuniao
}

func (repositorio *grupoRepositoryPostgres) ListarEventosAssociadosGrupo(grupoID string) []grupoService.AssociacaoGrupoEvento {
	repositorio.mutex.RLock()
	defer repositorio.mutex.RUnlock()
	return append([]grupoService.AssociacaoGrupoEvento(nil), repositorio.eventosAssociados[grupoID]...)
}

func (repositorio *grupoRepositoryPostgres) AssociarEventoGrupo(grupoID, eventoID string) grupoService.AssociacaoGrupoEvento {
	associacao := grupoService.AssociacaoGrupoEvento{ID: novoIdentificador("gev-"), GrupoID: grupoID, EventoID: eventoID}
	repositorio.mutex.Lock()
	repositorio.eventosAssociados[grupoID] = append(repositorio.eventosAssociados[grupoID], associacao)
	repositorio.mutex.Unlock()
	return associacao
}

func (repositorio *grupoRepositoryPostgres) ListarLeiturasAssociadasGrupo(grupoID string) []grupoService.AssociacaoGrupoLeitura {
	repositorio.mutex.RLock()
	defer repositorio.mutex.RUnlock()
	return append([]grupoService.AssociacaoGrupoLeitura(nil), repositorio.leiturasAssociadas[grupoID]...)
}

func (repositorio *grupoRepositoryPostgres) AssociarLeituraGrupo(grupoID, leituraID string) grupoService.AssociacaoGrupoLeitura {
	associacao := grupoService.AssociacaoGrupoLeitura{ID: novoIdentificador("grd-"), GrupoID: grupoID, LeituraID: leituraID}
	repositorio.mutex.Lock()
	repositorio.leiturasAssociadas[grupoID] = append(repositorio.leiturasAssociadas[grupoID], associacao)
	repositorio.mutex.Unlock()
	return associacao
}
