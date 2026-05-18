package service

import (
	"context"
	"errors"

	comum "campus_connect_api/internal/modulos/comum"
)

var ErrGrupoPrivado = errors.New("grupo privado")

func (servico *GrupoService) EntrarGrupo(contexto context.Context, grupoID, usuarioID string) error {
	vis, ok, err := servico.repositorio.ObterVisibilidadeGrupo(contexto, grupoID)
	if err != nil {
		return err
	}
	if !ok {
		return comum.ErrNaoEncontrado
	}
	if vis == "private" {
		return ErrGrupoPrivado
	}
	if err := servico.repositorio.InserirMembroGrupo(contexto, grupoID, usuarioID, "member"); err != nil {
		return err
	}
	return nil
}

func (servico *GrupoService) PedirEntradaGrupo(contexto context.Context, grupoID, usuarioID string) error {
	vis, ok, err := servico.repositorio.ObterVisibilidadeGrupo(contexto, grupoID)
	if err != nil {
		return err
	}
	if !ok {
		return comum.ErrNaoEncontrado
	}
	if vis != "private" {
		return servico.repositorio.InserirMembroGrupo(contexto, grupoID, usuarioID, "member")
	}
	return servico.repositorio.CriarPedidoEntradaGrupo(contexto, grupoID, usuarioID)
}

func (servico *GrupoService) ListarMembrosGrupo(contexto context.Context, grupoID string) ([]MembroGrupo, error) {
	return servico.repositorio.ListarMembrosGrupo(contexto, grupoID)
}
