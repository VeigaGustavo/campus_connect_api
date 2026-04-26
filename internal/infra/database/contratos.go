package database

import (
	"context"

	comunidadeModel "campus_connect_api/internal/modulos/comunidade/structs"
	discoverModel "campus_connect_api/internal/modulos/discover/structs"
	empresaModel "campus_connect_api/internal/modulos/empresa/structs"
	eventoModel "campus_connect_api/internal/modulos/evento/structs"
	grupoModel "campus_connect_api/internal/modulos/grupo/structs"
	leituraModel "campus_connect_api/internal/modulos/leitura/structs"
	perfilModel "campus_connect_api/internal/modulos/perfil/structs"
	projetoModel "campus_connect_api/internal/modulos/projeto/structs"
	universidadeModel "campus_connect_api/internal/modulos/universidade/structs"
	usuarioModel "campus_connect_api/internal/modulos/usuario/structs"
)

// Ports menores por agregado: o tipo Armazenamento continua sendo a união deles,
// facilitando testes e extração futura sem mudar as implementações Postgres/memória.

type PersistenciaUsuario interface {
	Autenticar(contexto context.Context, email, senha string) (*UsuarioInterno, error)
	CriarUsuario(contexto context.Context, nome, email, senha, perfilCodigo string) (*UsuarioInterno, error)
	CriarUsuarioComCadastro(contexto context.Context, req usuarioModel.RequisicaoCadastroUsuario) (*UsuarioInterno, error)
}

type PersistenciaOportunidade interface {
	ListarOportunidades(contexto context.Context) ([]empresaModel.Oportunidade, error)
	ObterOportunidade(contexto context.Context, id string) (empresaModel.Oportunidade, bool, error)
	InserirOportunidade(contexto context.Context, criadoPor string, corpo empresaModel.RequisicaoCriarOportunidade) (empresaModel.Oportunidade, error)
	AtualizarOportunidade(contexto context.Context, id, usuarioID, perfilCodigo string, corpo empresaModel.RequisicaoCriarOportunidade) (empresaModel.Oportunidade, error)
	RemoverOportunidade(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaEvento interface {
	ListarEventos(contexto context.Context) ([]eventoModel.EventoCampus, error)
	ObterEvento(contexto context.Context, id string) (eventoModel.EventoCampus, bool, error)
	InserirEvento(contexto context.Context, criadoPor string, corpo eventoModel.RequisicaoEvento) (eventoModel.EventoCampus, error)
	AtualizarEvento(contexto context.Context, id, usuarioID, perfilCodigo string, corpo eventoModel.RequisicaoEvento) (eventoModel.EventoCampus, error)
	RemoverEvento(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaGrupo interface {
	ListarGrupos(contexto context.Context) ([]grupoModel.GrupoEstudo, error)
	ObterGrupo(contexto context.Context, id string) (grupoModel.GrupoEstudo, bool, error)
	InserirGrupo(contexto context.Context, criadoPor string, corpo grupoModel.RequisicaoCriarGrupo) (grupoModel.GrupoEstudo, error)
	AtualizarGrupo(contexto context.Context, id, usuarioID, perfilCodigo string, corpo grupoModel.RequisicaoCriarGrupo) (grupoModel.GrupoEstudo, error)
	RemoverGrupo(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaComunidade interface {
	ListarComunidades(contexto context.Context) ([]comunidadeModel.Comunidade, error)
	ObterComunidade(contexto context.Context, id string) (comunidadeModel.Comunidade, bool, error)
	InserirComunidade(contexto context.Context, criadoPor string, corpo comunidadeModel.RequisicaoCriarComunidade) (comunidadeModel.Comunidade, error)
	AtualizarComunidade(contexto context.Context, id, usuarioID, perfilCodigo string, corpo comunidadeModel.RequisicaoCriarComunidade) (comunidadeModel.Comunidade, error)
	RemoverComunidade(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaAvisoUniversidade interface {
	ListarAvisosUniversidade(contexto context.Context) ([]universidadeModel.AvisoUniversidade, error)
	ObterAvisoUniversidade(contexto context.Context, id string) (universidadeModel.AvisoUniversidade, bool, error)
	InserirAvisoUniversidade(contexto context.Context, criadoPor string, corpo universidadeModel.RequisicaoCriarAvisoUniversidade) (universidadeModel.AvisoUniversidade, error)
	AtualizarAvisoUniversidade(contexto context.Context, id, usuarioID, perfilCodigo string, corpo universidadeModel.RequisicaoCriarAvisoUniversidade) (universidadeModel.AvisoUniversidade, error)
	RemoverAvisoUniversidade(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaLeitura interface {
	ListarLeituraSemanal(contexto context.Context) ([]leituraModel.ItemLeituraSemanal, error)
	ObterLeituraSemanal(contexto context.Context, id string) (leituraModel.ItemLeituraSemanal, bool, error)
	InserirLeituraSemanal(contexto context.Context, criadoPor string, corpo leituraModel.RequisicaoLeituraSemanal) (leituraModel.ItemLeituraSemanal, error)
	AtualizarLeituraSemanal(contexto context.Context, id, usuarioID, perfilCodigo string, corpo leituraModel.RequisicaoLeituraSemanal) (leituraModel.ItemLeituraSemanal, error)
	RemoverLeituraSemanal(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaProjeto interface {
	ListarProjetos(contexto context.Context) ([]projetoModel.Projeto, error)
	ObterProjeto(contexto context.Context, id string) (projetoModel.Projeto, bool, error)
	InserirProjeto(contexto context.Context, criadoPor string, corpo projetoModel.RequisicaoProjeto) (projetoModel.Projeto, error)
	AtualizarProjeto(contexto context.Context, id, usuarioID, perfilCodigo string, corpo projetoModel.RequisicaoProjeto) (projetoModel.Projeto, error)
	RemoverProjeto(contexto context.Context, id, usuarioID, perfilCodigo string) error
}

type PersistenciaFeedPerfil interface {
	FeedDescobrir(contexto context.Context, filtro string, gruposDoUsuario []string) ([]discoverModel.ItemDescobrir, error)
	PerfilUsuario(contexto context.Context, usuarioID string) (perfilModel.PerfilUsuario, error)
}

// Armazenamento é a união dos ports de persistência (implementação única em Postgres/memória).
type Armazenamento interface {
	PersistenciaUsuario
	PersistenciaOportunidade
	PersistenciaEvento
	PersistenciaGrupo
	PersistenciaComunidade
	PersistenciaAvisoUniversidade
	PersistenciaLeitura
	PersistenciaProjeto
	PersistenciaFeedPerfil
}
