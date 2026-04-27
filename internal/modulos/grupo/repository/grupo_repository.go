package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	comum "campus_connect_api/internal/modulos/comum"
	repositoryutil "campus_connect_api/internal/modulos/comum/repositoryutil"
	grupoService "campus_connect_api/internal/modulos/grupo/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type grupoRepositoryPostgres struct {
	pool *pgxpool.Pool

	mutex              sync.RWMutex
	chatGrupo          map[string][]grupoService.MensagemChatGrupo
	arquivosGrupo      map[string][]grupoService.ArquivoGrupo
	reunioesGrupo      map[string][]grupoService.ReuniaoGrupo
	eventosAssociados  map[string][]grupoService.AssociacaoGrupoEvento
	leiturasAssociadas map[string][]grupoService.AssociacaoGrupoLeitura
}

func NovoGrupoRepository(pool *pgxpool.Pool) grupoService.GrupoRepository {
	return &grupoRepositoryPostgres{
		pool:               pool,
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
	const sql = `SELECT id::text, titulo, field_of_study, description, level::text, member_count, schedule_label, criado_por::text FROM grupos_estudo ORDER BY criado_em DESC`
	rows, err := repositorio.pool.Query(contexto, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []grupoService.GrupoEstudo
	for rows.Next() {
		var g grupoService.GrupoEstudo
		var lvl string
		if err := rows.Scan(&g.Identificador, &g.Titulo, &g.AreaEstudo, &g.Descricao, &lvl, &g.TotalMembros, &g.RotuloHorario, &g.AutorID); err != nil {
			return nil, err
		}
		if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, g.AutorID, &g.Autor); err != nil {
			return nil, err
		}
		g.Nivel = grupoService.NivelGrupoEstudo(lvl)
		out = append(out, g)
	}
	return out, rows.Err()
}

func (repositorio *grupoRepositoryPostgres) InserirGrupo(contexto context.Context, criadoPor string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const ins = `INSERT INTO grupos_estudo (titulo, field_of_study, description, level, member_count, schedule_label, criado_por)
VALUES ($1,$2,$3,$4::varchar,$5,$6,$7::uuid) RETURNING id::text`
	var id string
	if err := tx.QueryRow(contexto, ins, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, 0, corpo.RotuloHorario, criadoPor).Scan(&id); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	if err := repositoryutil.InserirCartaoFeedTx(contexto, tx, comum.FeedKindGrupoEstudo, "dsc-"+id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel, id, corpo.EscopoPublicacao, corpo.IDGrupoPublicacao); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	if err := tx.Commit(contexto); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	g, ok, err := repositorio.obterGrupo(contexto, id)
	if err != nil || !ok {
		return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo")
	}
	return g, nil
}

func (repositorio *grupoRepositoryPostgres) AtualizarGrupo(contexto context.Context, id, usuarioID string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return repositorio.atualizarGrupoComPerfil(contexto, id, usuarioID, comum.PerfilPadrao, corpo)
}

func (repositorio *grupoRepositoryPostgres) AtualizarGrupoComoAdmin(contexto context.Context, id string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	return repositorio.atualizarGrupoComPerfil(contexto, id, "", comum.PerfilSistemaAdmin, corpo)
}

func (repositorio *grupoRepositoryPostgres) RemoverGrupo(contexto context.Context, id, usuarioID string) error {
	return repositorio.removerGrupoComPerfil(contexto, id, usuarioID, comum.PerfilPadrao)
}

func (repositorio *grupoRepositoryPostgres) RemoverGrupoComoAdmin(contexto context.Context, id string) error {
	return repositorio.removerGrupoComPerfil(contexto, id, "", comum.PerfilSistemaAdmin)
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

func (repositorio *grupoRepositoryPostgres) obterGrupo(contexto context.Context, id string) (grupoService.GrupoEstudo, bool, error) {
	const sql = `SELECT id::text, titulo, field_of_study, description, level::text, member_count, schedule_label, criado_por::text FROM grupos_estudo WHERE id=$1::uuid`
	var g grupoService.GrupoEstudo
	var lvl string
	err := repositorio.pool.QueryRow(contexto, sql, id).Scan(&g.Identificador, &g.Titulo, &g.AreaEstudo, &g.Descricao, &lvl, &g.TotalMembros, &g.RotuloHorario, &g.AutorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return grupoService.GrupoEstudo{}, false, nil
		}
		return grupoService.GrupoEstudo{}, false, err
	}
	if err := repositoryutil.CarregarPerfilPublicoAutor(contexto, repositorio.pool, g.AutorID, &g.Autor); err != nil {
		return grupoService.GrupoEstudo{}, false, err
	}
	g.Nivel = grupoService.NivelGrupoEstudo(lvl)
	return g, true, nil
}

func (repositorio *grupoRepositoryPostgres) atualizarGrupoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string, corpo grupoService.RequisicaoCriarGrupo) (grupoService.GrupoEstudo, error) {
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM grupos_estudo WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	const upd = `UPDATE grupos_estudo SET titulo=$2, field_of_study=$3, description=$4, level=$5::varchar, schedule_label=$6, atualizado_em=now() WHERE id=$1::uuid`
	ct, err := tx.Exec(contexto, upd, id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, corpo.Nivel, corpo.RotuloHorario)
	if err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	if ct.RowsAffected() == 0 {
		return grupoService.GrupoEstudo{}, comum.ErrNaoEncontrado
	}
	_, _ = tx.Exec(contexto, `UPDATE feed_cartoes SET titulo=$2, subtitle=$3, excerpt=$4, meta_primary=$5, meta_secondary=$6 WHERE kind='study_group' AND reference_id=$1`,
		id, corpo.Titulo, corpo.AreaEstudo, corpo.Descricao, "Nível", corpo.Nivel)
	if err := tx.Commit(contexto); err != nil {
		return grupoService.GrupoEstudo{}, err
	}
	g, ok, err := repositorio.obterGrupo(contexto, id)
	if err != nil || !ok {
		return grupoService.GrupoEstudo{}, errors.New("falha ao recarregar grupo")
	}
	return g, nil
}

func (repositorio *grupoRepositoryPostgres) removerGrupoComPerfil(contexto context.Context, id, usuarioID, perfilCodigo string) error {
	if err := repositoryutil.GarantirDonoOuAdmin(contexto, repositorio.pool, `SELECT criado_por::text FROM grupos_estudo WHERE id=$1::uuid`, id, usuarioID, perfilCodigo); err != nil {
		return err
	}
	tx, err := repositorio.pool.Begin(contexto)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(contexto) }()
	_, _ = tx.Exec(contexto, `DELETE FROM feed_cartoes WHERE kind='study_group' AND reference_id=$1`, id)
	ct, err := tx.Exec(contexto, `DELETE FROM grupos_estudo WHERE id=$1::uuid`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return comum.ErrNaoEncontrado
	}
	return tx.Commit(contexto)
}

