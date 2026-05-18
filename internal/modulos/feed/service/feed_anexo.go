package service

import (
	"context"
	"io"
	"strings"

	feedMedia "campus_connect_api/internal/modulos/feed/media"
)

type RespostaUploadAnexo struct {
	URL      string `json:"url"`
	Tipo     string `json:"type"`
	Nome     string `json:"name"`
	FileURL  string `json:"file_url"`
	MediaURL string `json:"media_url"`
	Filename string `json:"filename"`
}

func (servico *FeedService) EnviarAnexoFeed(_ context.Context, tipo string, nomeOriginal string, origem io.Reader) (RespostaUploadAnexo, error) {
	tipo = strings.ToLower(strings.TrimSpace(tipo))
	if tipo != "image" && tipo != "video" {
		return RespostaUploadAnexo{}, feedMedia.ErrTipoAnexo
	}
	caminho, err := feedMedia.SalvarAnexoFeed(tipo, nomeOriginal, origem)
	if err != nil {
		return RespostaUploadAnexo{}, err
	}
	nome := feedMedia.NomeExibicao(nomeOriginal, caminho)
	return RespostaUploadAnexo{
		URL:      caminho,
		Tipo:     tipo,
		Nome:     nome,
		FileURL:  caminho,
		MediaURL: caminho,
		Filename: nome,
	}, nil
}
