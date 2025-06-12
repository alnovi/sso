package http

import (
	"github.com/alnovi/gomon/server"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/provider"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/alnovi/sso/web"
)

func NewServer(p *provider.Provider) *server.HttpServer {
	mdwAuthSession := middleware.AuthBySession(p.Profile())
	mdwAdminToken := middleware.Token(p.OAuth(), p.Cookie(), p.Config().CAdmin.Id, p.Config().CAdmin.Secret)
	mdwAdminAuth := middleware.Auth(p.OAuth(), p.Cookie(), p.Config().CAdmin.Id, p.Config().CAdmin.Secret)
	mdwRoleAdmin := middleware.RoleWeight(entity.RoleAdminWeight)

	controllers := []server.HttpController{
		controller.NewProfileController(p.Profile(), p.Cookie(), mdwAuthSession),
		controller.NewAdminController(p.Admin(), p.Cookie(), mdwAdminToken),
		server.NewWrap("/oauth", []server.HttpController{
			oauth.NewCertsController(p.Certs()),
			oauth.NewAuthController(p.OAuth(), p.Cookie()),
			oauth.NewTokenController(p.OAuth()),
			oauth.NewPasswordController(p.OAuth()),
		}...),
		server.NewWrap("/api", []server.HttpController{
			api.NewClientController(p.StorageClients()),
			api.NewUserController(p.StorageUsers(), p.StorageRoles()),
			api.NewSessionController(p.StorageSessions()),
			api.NewStatsController(p.Stats()),
		}...).Use(mdwAdminAuth, mdwRoleAdmin),
	}

	s := server.NewHttpServer(
		server.WithHideBanner(),
		server.WithHidePort(),
		server.WithRender(server.NewHttpRenderFromFS(web.StaticFS, "out/html")),
		server.WithErrorHandler(controller.NewErrorController().Handle),
		server.WithValidator(p.Validator()),
		server.WithControllers(controllers...),
	)

	s.Pre(middleware.TrailingSlash())
	s.Use(middleware.Tracer())
	s.Use(middleware.RequestLogger(p.LoggerMod("http-request")))

	s.FileFS("/favicon.png/", "public/sso.png", web.StaticFS)
	s.StaticFS("/assets/*", echo.MustSubFS(web.StaticFS, "out/assets"))
	s.StaticFS("/public/*", echo.MustSubFS(web.StaticFS, "public"))
	s.GET("/swagger/*", echoSwagger.WrapHandler)

	p.Closer().Add(s.Shutdown)

	return s
}
