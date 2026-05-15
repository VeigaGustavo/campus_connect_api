package media

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxBytesImagemFeed = 8 << 20  // 8 MiB
	MaxBytesVideoFeed  = 80 << 20 // 80 MiB
	subpastaFeed       = "feed"
)

var (
	ErrArquivoGrande   = errors.New("arquivo excede o tamanho maximo")
	ErrFormatoInvalido = errors.New("formato de arquivo invalido")
	ErrTipoAnexo       = errors.New("tipo de anexo invalido")
	ErrArquivoVazio    = errors.New("arquivo vazio")
)

// SalvarAnexoFeed grava imagem ou video em uploads/feed e devolve caminho /uploads/feed/...
func SalvarAnexoFeed(tipo string, nomeOriginal string, origem io.Reader) (caminhoRelativo string, err error) {
	tipo = strings.ToLower(strings.TrimSpace(tipo))
	switch tipo {
	case "image", "video":
	default:
		return "", ErrTipoAnexo
	}
	limite := MaxBytesImagemFeed
	if tipo == "video" {
		limite = MaxBytesVideoFeed
	}
	dados, err := io.ReadAll(io.LimitReader(origem, int64(limite)+1))
	if err != nil {
		return "", err
	}
	if len(dados) == 0 {
		return "", ErrArquivoVazio
	}
	if len(dados) > limite {
		return "", ErrArquivoGrande
	}
	ext, err := extensaoPorTipo(tipo, dados, nomeOriginal)
	if err != nil {
		return "", err
	}
	nomeArquivo, err := nomeAleatorio(ext)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(ResolverDirUploads(), subpastaFeed)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	caminho := filepath.Join(dir, nomeArquivo)
	if err := os.WriteFile(caminho, dados, 0o644); err != nil {
		return "", err
	}
	return "/uploads/" + subpastaFeed + "/" + nomeArquivo, nil
}

func ResolverDirUploads() string {
	if v := strings.TrimSpace(os.Getenv("UPLOAD_DIR")); v != "" {
		return v
	}
	return "./data/uploads"
}

func extensaoPorTipo(tipo string, dados []byte, nomeOriginal string) (string, error) {
	if tipo == "video" {
		if ext := extensaoVideo(dados, nomeOriginal); ext != "" {
			return ext, nil
		}
		return "", ErrFormatoInvalido
	}
	if ext := extensaoImagem(dados); ext != "" {
		return ext, nil
	}
	return "", ErrFormatoInvalido
}

func extensaoImagem(dados []byte) string {
	if len(dados) >= 3 && dados[0] == 0xFF && dados[1] == 0xD8 {
		return ".jpg"
	}
	if len(dados) >= 8 && string(dados[0:8]) == "\x89PNG\r\n\x1a\n" {
		return ".png"
	}
	if len(dados) >= 6 && (string(dados[0:6]) == "GIF87a" || string(dados[0:6]) == "GIF89a") {
		return ".gif"
	}
	if len(dados) >= 12 && string(dados[0:4]) == "RIFF" && string(dados[8:12]) == "WEBP" {
		return ".webp"
	}
	return ""
}

func extensaoVideo(dados []byte, nomeOriginal string) string {
	if len(dados) >= 12 {
		if string(dados[4:8]) == "ftyp" {
			return ".mp4"
		}
		if len(dados) >= 4 && string(dados[0:4]) == "\x1a\x45\xdf\xa3" {
			return ".webm"
		}
	}
	nome := strings.ToLower(filepath.Ext(nomeOriginal))
	switch nome {
	case ".mp4", ".webm", ".mov", ".m4v":
		return nome
	}
	return ""
}

func nomeAleatorio(ext string) (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", hex.EncodeToString(buf[:]), ext), nil
}

// NomeExibicao devolve nome seguro para resposta (original ou derivado do path).
func NomeExibicao(nomeOriginal, caminhoRelativo string) string {
	nome := strings.TrimSpace(filepath.Base(nomeOriginal))
	if nome != "" && nome != "." {
		return nome
	}
	return filepath.Base(caminhoRelativo)
}
