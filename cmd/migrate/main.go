package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"campus_connect_api/internal/banco"

	"github.com/joho/godotenv"
)

func main() {
	carregarEnvLocal()

	contexto := context.Background()
	urlBanco := resolverDatabaseURL()
	if urlBanco == "" {
		log.Fatal("DATABASE_URL ou DATABASE_URL_FILE é obrigatório")
	}

	pool, err := banco.NovoPool(contexto, urlBanco)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	if err := banco.AplicarMigracoes(contexto, pool); err != nil {
		log.Fatalf("migracoes: %v", err)
	}
	log.Println("migracoes aplicadas com sucesso")
}

func carregarEnvLocal() {
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) != "" || strings.TrimSpace(os.Getenv("DATABASE_URL_FILE")) != "" {
		return
	}
	candidatos := []string{".env.local", ".env"}
	if _, arquivoFonte, _, ok := runtime.Caller(0); ok {
		dirProjeto := filepath.Dir(filepath.Dir(filepath.Dir(arquivoFonte)))
		candidatos = append(candidatos,
			filepath.Join(dirProjeto, ".env.local"),
			filepath.Join(dirProjeto, ".env"),
		)
	}
	for _, caminho := range candidatos {
		_ = godotenv.Load(caminho)
	}
}

func resolverDatabaseURL() string {
	if url := strings.TrimSpace(os.Getenv("DATABASE_URL")); url != "" {
		return url
	}
	arquivo := strings.TrimSpace(os.Getenv("DATABASE_URL_FILE"))
	if arquivo == "" {
		return ""
	}
	dados, err := os.ReadFile(filepath.Clean(arquivo))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		log.Printf("erro ao ler DATABASE_URL_FILE: %v", err)
		return ""
	}
	return strings.TrimSpace(string(dados))
}
