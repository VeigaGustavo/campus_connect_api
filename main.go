package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"campus_connect_api/internal/app"
	"campus_connect_api/internal/banco"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	contexto := context.Background()
	urlBancoDados := resolverDatabaseURL()
	if urlBancoDados == "" {
		log.Fatal("DATABASE_URL ou DATABASE_URL_FILE é obrigatório; configure o PostgreSQL e rode as migrações em db/init")
	}
	pool, err := banco.NovoPool(contexto, urlBancoDados)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()
	log.Println("persistência: PostgreSQL (DATABASE_URL)")

	engine := app.NewGinEngine(pool)
	endereco := resolverEnderecoEscuta()

	srv := &http.Server{
		Addr:              endereco,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("campus_connect_api escutando em %s (API em /api/)", endereco)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func resolverEnderecoEscuta() string {
	if definido := os.Getenv("LISTEN_ADDR"); definido != "" {
		return definido
	}
	if porta := os.Getenv("PORT"); porta != "" {
		return ":" + porta
	}
	return ":8080"
}

func resolverDatabaseURL() string {
	if url := strings.TrimSpace(os.Getenv("DATABASE_URL")); url != "" {
		return url
	}

	arquivoSecret := strings.TrimSpace(os.Getenv("DATABASE_URL_FILE"))
	if arquivoSecret == "" {
		return ""
	}

	dados, err := os.ReadFile(filepath.Clean(arquivoSecret))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("DATABASE_URL_FILE não encontrado: %s", arquivoSecret)
			return ""
		}
		log.Printf("erro ao ler DATABASE_URL_FILE: %v", err)
		return ""
	}

	return strings.TrimSpace(string(dados))
}
