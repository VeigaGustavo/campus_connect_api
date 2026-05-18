package media

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SalvarJPEG(dirBase, subpasta, nomeArquivo string, conteudo []byte) (string, error) {
	dir := filepath.Join(dirBase, subpasta)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	caminho := filepath.Join(dir, nomeArquivo)
	if err := os.WriteFile(caminho, conteudo, 0o644); err != nil {
		return "", err
	}
	return "/uploads/" + subpasta + "/" + nomeArquivo, nil
}

func ResolverDirUploads() string {
	if v := strings.TrimSpace(os.Getenv("UPLOAD_DIR")); v != "" {
		return v
	}
	return "./data/uploads"
}

func AbsolutizarURL(caminho string) string {
	caminho = strings.TrimSpace(caminho)
	if caminho == "" {
		return ""
	}
	if strings.HasPrefix(caminho, "http://") || strings.HasPrefix(caminho, "https://") {
		return caminho
	}
	return URLPublica(caminho)
}

func URLPublica(caminhoRelativo string) string {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")), "/")
	if base == "" {
		return caminhoRelativo
	}
	if !strings.HasPrefix(caminhoRelativo, "/") {
		caminhoRelativo = "/" + caminhoRelativo
	}
	return base + caminhoRelativo
}

func NomeArquivoPerfil(usuarioID, tipo string) string {
	id := strings.ReplaceAll(usuarioID, "-", "")
	return fmt.Sprintf("%s_%s.jpg", tipo, id)
}
