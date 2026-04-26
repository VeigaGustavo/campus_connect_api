package service

import model "campus_connect_api/internal/modulos/empresa/structs"

type ModalidadeLocalTrabalho = model.ModalidadeLocalTrabalho

const (
	TrabalhoRemoto     = model.TrabalhoRemoto
	TrabalhoHibrido    = model.TrabalhoHibrido
	TrabalhoPresencial = model.TrabalhoPresencial
)

type Oportunidade = model.Oportunidade
type RequisicaoCriarOportunidade = model.RequisicaoCriarOportunidade
type CandidatoOportunidade = model.CandidatoOportunidade
