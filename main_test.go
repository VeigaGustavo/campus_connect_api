package main

import "testing"

func TestResolverEnderecoEscuta(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		t.Setenv("LISTEN_ADDR", "")
		t.Setenv("PORT", "")

		if got := resolverEnderecoEscuta(); got != ":8080" {
			t.Fatalf("esperado :8080, obtido %q", got)
		}
	})

	t.Run("porta", func(t *testing.T) {
		t.Setenv("LISTEN_ADDR", "")
		t.Setenv("PORT", "9090")

		if got := resolverEnderecoEscuta(); got != ":9090" {
			t.Fatalf("esperado :9090, obtido %q", got)
		}
	})

	t.Run("listen_addr tem prioridade", func(t *testing.T) {
		t.Setenv("LISTEN_ADDR", "127.0.0.1:7070")
		t.Setenv("PORT", "9090")

		if got := resolverEnderecoEscuta(); got != "127.0.0.1:7070" {
			t.Fatalf("esperado 127.0.0.1:7070, obtido %q", got)
		}
	})
}

