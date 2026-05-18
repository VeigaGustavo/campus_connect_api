package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"syscall"

	"campus_connect_api/internal/app"
	"campus_connect_api/internal/banco"
	usuarioRepository "campus_connect_api/internal/modulos/usuario/repository"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))
	carregarEnvLocal()

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
	if err := banco.AplicarMigracoesEssenciais(contexto, pool); err != nil {
		log.Fatalf("migracoes: %v", err)
	}
	if n, err := usuarioRepository.NovoUsuarioRepository(pool).RepararComunidadesSemGrupo(contexto); err != nil {
		log.Printf("aviso: reparar comunidades sem grupo: %v", err)
	} else if n > 0 {
		log.Printf("reparadas %d conta(s) comunidade sem grupo (atlética/CA)", n)
	}
	log.Println("persistência: PostgreSQL (DATABASE_URL)")

	engine := app.NewGinEngine(pool)
	endereco := resolverEnderecoEscuta()

	listener, err := escutarHTTP(endereco)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	srv := &http.Server{
		Addr:              listener.Addr().String(),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("campus_connect_api escutando em %s (API em /api/)", listener.Addr())
	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func carregarEnvLocal() {
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) != "" || strings.TrimSpace(os.Getenv("DATABASE_URL_FILE")) != "" {
		return
	}

	candidatos := []string{".env.local", ".env"}
	if _, arquivoFonte, _, ok := runtime.Caller(0); ok {
		dirFonte := filepath.Dir(arquivoFonte)
		candidatos = append(candidatos,
			filepath.Join(dirFonte, ".env.local"),
			filepath.Join(dirFonte, ".env"),
		)
	}

	visitados := make(map[string]struct{}, len(candidatos))
	for _, caminho := range candidatos {
		if caminho == "" {
			continue
		}
		abs, err := filepath.Abs(caminho)
		if err != nil {
			abs = caminho
		}
		if _, existe := visitados[abs]; existe {
			continue
		}
		visitados[abs] = struct{}{}

		if err := godotenv.Load(abs); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("aviso: erro ao carregar arquivo de ambiente %s: %v", abs, err)
		}
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

func escutarHTTP(endereco string) (net.Listener, error) {
	listener, err := net.Listen("tcp", endereco)
	if err == nil {
		return listener, nil
	}

	if endereco == ":8080" && errors.Is(err, syscall.EADDRINUSE) {
		log.Printf("porta 8080 ocupada; usando uma porta livre automática para esta execucao")
		return net.Listen("tcp", ":0")
	}

	return nil, err
}

func resolverDatabaseURL() string {
	if url := normalizarValorEnv(os.Getenv("DATABASE_URL")); url != "" {
		return url
	}

	arquivoSecret := normalizarValorEnv(os.Getenv("DATABASE_URL_FILE"))
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

func normalizarValorEnv(valor string) string {
	valor = strings.TrimSpace(valor)
	if strings.HasPrefix(valor, "=") {
		valor = strings.TrimSpace(strings.TrimPrefix(valor, "="))
	}
	return valor
}

