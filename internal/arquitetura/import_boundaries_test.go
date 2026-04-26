package arquitetura

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	modulosPrefixo      = `"campus_connect_api/internal/modulos/`
	infraDatabaseImport = `"campus_connect_api/internal/infra/database"`
)

func TestServiceNaoImportaInfraDatabase(t *testing.T) {
	validarArquivos(t, "service", func(caminho, conteudo string) {
		if strings.Contains(conteudo, infraDatabaseImport) {
			t.Errorf("service nao pode importar infra/database: %s", caminho)
		}
	})
}

func TestServiceNaoImportaRepositoryNemHandlerDeModulo(t *testing.T) {
	validarArquivos(t, "service", func(caminho, conteudo string) {
		if strings.Contains(conteudo, "/repository") || strings.Contains(conteudo, "/handler") {
			t.Errorf("service nao pode importar repository/handler: %s", caminho)
		}
	})
}

func TestRepositoryNaoImportaHandler(t *testing.T) {
	validarArquivos(t, "repository", func(caminho, conteudo string) {
		if strings.Contains(conteudo, modulosPrefixo) && strings.Contains(conteudo, "/handler") {
			t.Errorf("repository nao pode importar handler: %s", caminho)
		}
	})
}

func validarArquivos(t *testing.T, camada string, validar func(caminho, conteudo string)) {
	t.Helper()

	baseModulos := filepath.Join("..", "modulos")
	err := filepath.WalkDir(baseModulos, func(caminho string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(caminho, ".go") {
			return nil
		}
		alvoCamada := string(filepath.Separator) + camada + string(filepath.Separator)
		if !strings.Contains(caminho, alvoCamada) {
			return nil
		}

		bruto, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		validar(caminho, string(bruto))
		return nil
	})

	if err != nil {
		t.Fatalf("erro ao validar boundaries da camada %s: %v", camada, err)
	}
}
