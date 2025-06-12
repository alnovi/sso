package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alnovi/gomon/server"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/provider"
	trHttp "github.com/alnovi/sso/internal/transport/http"
)

type App struct {
	Provider   *provider.Provider
	HttpServer *server.HttpServer
}

func NewApp(cfg *config.Config) *App {
	p := provider.New(cfg)

	defer func() {
		if err := recover(); err != nil {
			p.LoggerMod("app-server").Error(fmt.Sprintf("failed init app: %s", err.(error).Error()))
			os.Exit(1)
		}
	}()

	docs.SwaggerInfo.Version = config.Version
	docs.SwaggerInfo.Host = p.Config().App.Host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	p.Tracer()
	p.MigrationUp()

	return &App{Provider: p, HttpServer: trHttp.NewServer(p)}
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

	app.Provider.LoggerMod("app-server").Info("server started",
		slog.String("http", app.Provider.Config().Http.Addr()),
		slog.String("version", config.Version),
	)

	<-ctx.Done()
}
