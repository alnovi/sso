package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/provider"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/alnovi/sso/pkg/server"
	"github.com/alnovi/sso/web"
)

type App struct {
	Provider    *provider.Provider
	Controllers []server.HttpController
	HttpServer  *server.HttpServer
}

func NewApp(cfg *config.Config) *App {
	app := &App{Provider: provider.New(cfg)}

	defer func() {
		if err := recover(); err != nil {
			app.Provider.LoggerMod("app-server").Error(fmt.Sprintf("failed init app: %s", err.(error).Error()))
			os.Exit(1)
		}
	}()

	app.Provider.MigrationUp()
	app.initControllers()
	app.initHTTPServer()
	app.initSwag()

	return app
}

func (app *App) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := app.Provider.Closer().Close(); err != nil {
			app.Provider.LoggerMod("closer").Error(err.Error())
		}

		if err := recover(); err != nil {
			app.Provider.LoggerMod("app-server").Error(err.(error).Error())
			os.Exit(1)
		}

		cancel()
	}()

	go func() {
		app.Provider.Scheduler().Start()
	}()

	go func() {
		err := app.HttpServer.Start(app.Provider.Config().Http.Host, app.Provider.Config().Http.Port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Provider.LoggerMod("http-server").Error(err.(error).Error())
			cancel()
		}
	}()

	app.Provider.LoggerMod("http-server").Info("server started", "port", app.Provider.Config().Http.Port)

	<-ctx.Done()
}

func (app *App) initControllers() {
	mdwAuthSession := middleware.AuthBySession(app.Provider.Profile())
	mdwAdminToken := middleware.Token(app.Provider.OAuth(), app.Provider.Cookie(), app.Provider.Config().CAdmin.Id, app.Provider.Config().CAdmin.Secret)
	mdwAdminAuth := middleware.Auth(app.Provider.OAuth(), app.Provider.Cookie(), app.Provider.Config().CAdmin.Id, app.Provider.Config().CAdmin.Secret)
	mdwRoleAdmin := middleware.RoleWeight(entity.RoleAdminWeight)

	app.Controllers = []server.HttpController{
		controller.NewProfileController(app.Provider.Profile(), app.Provider.Cookie(), mdwAuthSession),
		controller.NewAdminController(app.Provider.Admin(), app.Provider.Cookie(), mdwAdminToken),
		server.NewWrap("/oauth", []server.HttpController{
			oauth.NewAuthController(app.Provider.OAuth(), app.Provider.Cookie()),
			oauth.NewTokenController(app.Provider.OAuth()),
			oauth.NewPasswordController(app.Provider.OAuth()),
		}...),
		server.NewWrap("/api", []server.HttpController{
			api.NewClientController(app.Provider.StorageClients()),
			api.NewUserController(app.Provider.StorageUsers(), app.Provider.StorageRoles()),
			api.NewSessionController(app.Provider.StorageSessions()),
			api.NewStatsController(app.Provider.Stats()),
		}...).Use(mdwAdminAuth, mdwRoleAdmin),
	}
}

func (app *App) initHTTPServer() {
	app.HttpServer = server.NewHttpServer(
		server.WithHideBanner(),
		server.WithHidePort(),
		server.WithRender(server.NewHttpRenderFromFS(web.StaticFS, "out/html")),
		server.WithErrorHandler(controller.NewErrorController().Handle),
		server.WithValidator(app.Provider.Validator()),
		server.WithControllers(app.Controllers...),
		server.WithCors(app.Provider.Config().IsDevelopment()),
	)

	app.HttpServer.Pre(middleware.TrailingSlash())
	app.HttpServer.Use(middleware.RequestLogger(app.Provider.LoggerMod("http-request")))

	app.HttpServer.FileFS("/favicon.png/", "public/sso.png", web.StaticFS)
	app.HttpServer.StaticFS("/assets/*", echo.MustSubFS(web.StaticFS, "out/assets"))
	app.HttpServer.StaticFS("/public/*", echo.MustSubFS(web.StaticFS, "public"))
	app.HttpServer.GET("/swagger/*", echoSwagger.WrapHandler)

	app.Provider.Closer().Add(app.HttpServer.Shutdown)
}

func (app *App) initSwag() {
	docs.SwaggerInfo.Version = app.Provider.Config().App.Version
	docs.SwaggerInfo.Host = app.Provider.Config().App.Host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
