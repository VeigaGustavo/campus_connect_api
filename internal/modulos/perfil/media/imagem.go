package media

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"

	"golang.org/x/image/draw"
)

const (
	MaxBytesEntrada = 5 << 20 // 5 MiB
	qualidadeJPEG   = 82
	larguraAvatar   = 256
	alturaAvatar    = 256
	larguraCapa     = 1280
	alturaCapa      = 420
)

var (
	ErrArquivoGrande   = errors.New("arquivo de imagem excede o tamanho maximo")
	ErrFormatoInvalido = errors.New("formato de imagem nao suportado")
	ErrImagemVazia     = errors.New("imagem vazia ou invalida")
)

func ProcessarAvatar(origem io.Reader) ([]byte, error) {
	img, err := decodificarLimitado(origem)
	if err != nil {
		return nil, err
	}
	alvo := image.NewRGBA(image.Rect(0, 0, larguraAvatar, alturaAvatar))
	draw.CatmullRom.Scale(alvo, alvo.Bounds(), img, img.Bounds(), draw.Over, nil)
	return codificarJPEG(alvo)
}

func ProcessarCapa(origem io.Reader) ([]byte, error) {
	img, err := decodificarLimitado(origem)
	if err != nil {
		return nil, err
	}
	alvo := image.NewRGBA(image.Rect(0, 0, larguraCapa, alturaCapa))
	draw.CatmullRom.Scale(alvo, alvo.Bounds(), img, img.Bounds(), draw.Over, nil)
	return codificarJPEG(alvo)
}

func decodificarLimitado(origem io.Reader) (image.Image, error) {
	dados, err := io.ReadAll(io.LimitReader(origem, MaxBytesEntrada+1))
	if err != nil {
		return nil, err
	}
	if len(dados) == 0 {
		return nil, ErrImagemVazia
	}
	if len(dados) > MaxBytesEntrada {
		return nil, ErrArquivoGrande
	}
	img, _, err := image.Decode(bytes.NewReader(dados))
	if err != nil {
		return nil, ErrFormatoInvalido
	}
	return img, nil
}

func codificarJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: qualidadeJPEG}); err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		return nil, ErrImagemVazia
	}
	return buf.Bytes(), nil
}
