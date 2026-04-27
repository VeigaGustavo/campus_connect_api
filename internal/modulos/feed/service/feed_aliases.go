package service

import (
	comum "campus_connect_api/internal/modulos/comum"
	model "campus_connect_api/internal/modulos/feed/structs"
)

type CategoriaItemFeed = model.CategoriaItemFeed

const (
	CategoriaEstagio     = model.CategoriaEstagio
	CategoriaEvento      = model.CategoriaEvento
	CategoriaGrupoEstudo = model.CategoriaGrupoEstudo
	CategoriaProjeto     = model.CategoriaProjeto
	CategoriaAviso       = model.CategoriaAviso
	CategoriaLeitura     = model.CategoriaLeitura
	CategoriaPost        = model.CategoriaPost
)

type ItemFeed = model.ItemFeed
type RespostaFeed = model.RespostaFeed
type AnexoPost = model.AnexoPost
type PerfilPublicoAutor = comum.PerfilPublicoAutor
type RequisicaoCriarPost = model.RequisicaoCriarPost
type ComentarioPost = model.ComentarioPost
type RequisicaoCriarComentario = model.RequisicaoCriarComentario
type RequisicaoReacao = model.RequisicaoReacao
type RequisicaoSalvarPost = model.RequisicaoSalvarPost
type PostFeedDetalhe = model.PostFeedDetalhe
