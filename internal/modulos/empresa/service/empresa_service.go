package service

import (
	"context"
)

type EmpresaService struct {
	repositorio EmpresaRepository
}

func NovoEmpresaService(repositorio EmpresaRepository) *EmpresaService {
	return &EmpresaService{repositorio: repositorio}
}

func (servico *EmpresaService) ListarOportunidades(contexto context.Context) ([]Oportunidade, error) {
	return servico.repositorio.ListarOportunidades(contexto)
}

func (servico *EmpresaService) ObterOportunidadePorID(contexto context.Context, id string) (Oportunidade, bool, error) {
	return servico.repositorio.ObterOportunidade(contexto, id)
}

func (servico *EmpresaService) CriarOportunidade(contexto context.Context, criadoPor string, corpo RequisicaoCriarOportunidade) (Oportunidade, error) {
	return servico.repositorio.InserirOportunidade(contexto, criadoPor, corpo)
}

func (servico *EmpresaService) AtualizarOportunidade(contexto context.Context, id, usuarioID, perfil string, corpo RequisicaoCriarOportunidade) (Oportunidade, error) {
	return servico.repositorio.AtualizarOportunidade(contexto, id, usuarioID, perfil, corpo)
}

func (servico *EmpresaService) RemoverOportunidade(contexto context.Context, id, usuarioID, perfil string) error {
	return servico.repositorio.RemoverOportunidade(contexto, id, usuarioID, perfil)
}

func (servico *EmpresaService) ListarCandidatosOportunidade(oportunidadeID string) []CandidatoOportunidade {
	_ = oportunidadeID
	return []CandidatoOportunidade{
		{
			UsuarioID: "usr-pad-1", Nome: "Joao Oliveira", Email: "joao@campusconnect.com",
			ResumoPerfil: "Aluno de ADS com projetos em Go e Flutter.", URLCurriculo: "https://example.com/cv/joao.pdf",
			StatusCandidatura: "submitted",
		},
		{
			UsuarioID: "usr-pad-2", Nome: "Maria Lima", Email: "maria@campusconnect.com",
			ResumoPerfil: "Experiência em eventos acadêmicos e monitoria.", URLCurriculo: "https://example.com/cv/maria.pdf",
			StatusCandidatura: "in_review",
		},
	}
}
