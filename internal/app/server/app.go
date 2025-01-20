package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/provider"
	"github.com/alnovi/sso/internal/transaport/http/controller"
	"github.com/alnovi/sso/internal/transaport/http/middleware"
	"github.com/alnovi/sso/pkg/server"
)

type App struct {
	Provider    *provider.Provider
	HttpServer  *server.HttpServer
	Controllers []server.HttpController
}

func NewApp(cfg *config.Config) *App {
	_ = os.Setenv("TZ", "UTC")

	app := &App{Provider: provider.NewProvider(cfg)}

	defer func() {
		if err := recover(); err != nil {
			app.Provider.LoggerModule("app-server").Error(err.(error).Error())
			os.Exit(1)
		}
	}()

	app.Provider.MigrationUp(context.Background())
	app.initControllers()
	app.initHttpServer()
	app.initSwagger()

	return app
}

func (app *App) Start(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := app.Provider.Closer().Close(); err != nil {
			app.Provider.LoggerModule("closer").Error(err.Error())
		}

		if err := recover(); err != nil {
			app.Provider.LoggerModule("app-server").Error(err.(error).Error())
			os.Exit(1)
		}

		cancel()
	}()

	go func() {
		err := app.HttpServer.Start(app.Provider.Config().Http.Host, app.Provider.Config().Http.Port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Provider.LoggerModule("app-server").Error(err.(error).Error())
			cancel()
		}
	}()

	app.Provider.LoggerModule("app-server").Info(fmt.Sprintf("server started, listening on port: %s", app.Provider.Config().Http.Port))

	<-ctx.Done()
}

func (app *App) initControllers() {
	app.Controllers = []server.HttpController{}
}

func (app *App) initHttpServer() {
	app.HttpServer = server.NewHttpServer(
		server.WithHideBanner(),
		server.WithHidePort(),
		server.WithErrorHandler(controller.NewErrorHandler().Handle),
		server.WithValidator(app.Provider.Validator()),
		server.WithControllers(app.Controllers...),
		server.WithCors([]string{"http://127.0.0.1:8080"}, nil),
	)

	app.HttpServer.Pre(middleware.TrailingSlash())
	app.HttpServer.Use(middleware.RequestLogger(app.Provider.LoggerModule("http-request")))

	app.HttpServer.GET("/swagger/*", echoSwagger.WrapHandler)

	app.Provider.Closer().Add(app.HttpServer.Shutdown)
}

func (app *App) initSwagger() {
	docs.SwaggerInfo.Version = app.Provider.Config().App.Version
	docs.SwaggerInfo.Host = app.Provider.Config().App.Host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
