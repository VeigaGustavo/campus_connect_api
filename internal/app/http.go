package app

import (
	comunidadeHandler "campus_connect_api/internal/modulos/comunidade/handler"
	comunidadeRepository "campus_connect_api/internal/modulos/comunidade/repository"
	comunidadeService "campus_connect_api/internal/modulos/comunidade/service"
	empresaHandler "campus_connect_api/internal/modulos/empresa/handler"
	empresaRepository "campus_connect_api/internal/modulos/empresa/repository"
	empresaService "campus_connect_api/internal/modulos/empresa/service"
	eventoHandler "campus_connect_api/internal/modulos/evento/handler"
	eventoRepository "campus_connect_api/internal/modulos/evento/repository"
	eventoService "campus_connect_api/internal/modulos/evento/service"
	feedHandler "campus_connect_api/internal/modulos/feed/handler"
	feedRepository "campus_connect_api/internal/modulos/feed/repository"
	feedService "campus_connect_api/internal/modulos/feed/service"
	grupoHandler "campus_connect_api/internal/modulos/grupo/handler"
	grupoRepository "campus_connect_api/internal/modulos/grupo/repository"
	grupoService "campus_connect_api/internal/modulos/grupo/service"
	leituraHandler "campus_connect_api/internal/modulos/leitura/handler"
	leituraRepository "campus_connect_api/internal/modulos/leitura/repository"
	leituraService "campus_connect_api/internal/modulos/leitura/service"
	perfilHandler "campus_connect_api/internal/modulos/perfil/handler"
	"campus_connect_api/internal/modulos/perfil/media"
	perfilRepository "campus_connect_api/internal/modulos/perfil/repository"
	perfilService "campus_connect_api/internal/modulos/perfil/service"
	projetoHandler "campus_connect_api/internal/modulos/projeto/handler"
	projetoRepository "campus_connect_api/internal/modulos/projeto/repository"
	projetoService "campus_connect_api/internal/modulos/projeto/service"
	segurancaHandler "campus_connect_api/internal/modulos/seguranca/handler"
	segurancaRepository "campus_connect_api/internal/modulos/seguranca/repository"
	segurancaService "campus_connect_api/internal/modulos/seguranca/service"
	universidadeHandler "campus_connect_api/internal/modulos/universidade/handler"
	universidadeRepository "campus_connect_api/internal/modulos/universidade/repository"
	universidadeService "campus_connect_api/internal/modulos/universidade/service"
	usuarioHandler "campus_connect_api/internal/modulos/usuario/handler"
	usuarioRepository "campus_connect_api/internal/modulos/usuario/repository"
	usuarioService "campus_connect_api/internal/modulos/usuario/service"
	"campus_connect_api/internal/respostas"
	"campus_connect_api/internal/versao"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewGinEngine(pool *pgxpool.Pool) *gin.Engine {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	servicoFeed := feedService.NovoFeedService(feedRepository.NovoFeedRepository(pool))
	servicoEmpresa := empresaService.NovoEmpresaService(empresaRepository.NovoEmpresaRepository(pool))
	servicoEvento := eventoService.NovoEventoService(eventoRepository.NovoEventoRepository(pool))
	servicoGrupo := grupoService.NovoGrupoService(grupoRepository.NovoGrupoRepository(pool))
	servicoComunidade := comunidadeService.NovoComunidadeService(comunidadeRepository.NovoComunidadeRepository(pool))
	servicoUniversidade := universidadeService.NovoUniversidadeService(universidadeRepository.NovoUniversidadeRepository(pool))
	servicoLeitura := leituraService.NovoLeituraService(leituraRepository.NovoLeituraRepository(pool))
	servicoProjeto := projetoService.NovoProjetoService(projetoRepository.NovoProjetoRepository(pool))
	servicoPerfil := perfilService.NovoPerfilService(perfilRepository.NovoPerfilRepository(pool))
	servicoSeguranca := segurancaService.NovoSegurancaService(segurancaRepository.NovoSegurancaRepository(pool))
	servicoUsuario := usuarioService.NovoUsuarioService(usuarioRepository.NovoUsuarioRepository(pool))

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(respostas.GinRequestID())
	engine.Use(respostas.GinAPIRevision())
	engine.Use(respostas.GinCORS())
	engine.Use(respostas.GinAceitarJSON())
	engine.GET("/health", func(contexto *gin.Context) {
		contexto.JSON(200, map[string]any{
			"status":       "ok",
			"api_revision": versao.Revisao,
			"features":     versao.Features,
		})
	})
	engine.Static("/uploads", media.ResolverDirUploads())

	api := engine.Group("/api")
	segurancaHandler.NovoSegurancaHTTPHandler(servicoSeguranca).RegistrarRotasGIN(api)
	usuarioHandler.NovoUsuarioHTTPHandler(servicoUsuario).RegistrarRotasGIN(api)
	feedHandler.NovoFeedHTTPHandler(servicoFeed).RegistrarRotasGIN(api)
	empresaHandler.NovoEmpresaHTTPHandler(servicoEmpresa).RegistrarRotasGIN(api)
	eventoHandler.NovoEventoHTTPHandler(servicoEvento).RegistrarRotasGIN(api)
	grupoHandler.NovoGrupoHTTPHandler(servicoGrupo).RegistrarRotasGIN(api)
	comunidadeHandler.NovoComunidadeHTTPHandler(servicoComunidade).RegistrarRotasGIN(api)
	universidadeHandler.NovoUniversidadeHTTPHandler(servicoUniversidade).RegistrarRotasGIN(api)
	leituraHandler.NovoLeituraHTTPHandler(servicoLeitura).RegistrarRotasGIN(api)
	projetoHandler.NovoProjetoHTTPHandler(servicoProjeto).RegistrarRotasGIN(api)
	perfilHTTP := perfilHandler.NovoPerfilHTTPHandler(servicoPerfil)
	perfilHTTP.RegistrarRotasUploadGIN(api)
	perfilHTTP.RegistrarRotasGIN(api)

	return engine
}
