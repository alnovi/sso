package server

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alnovi/sso/internal/config"
	"github.com/alnovi/sso/internal/transport/http/handler"
	"github.com/alnovi/sso/internal/transport/http/render"
	"github.com/alnovi/sso/pkg/server"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/alnovi/sso/web"
)

type App struct {
	*Provider
	Server *server.Server
}

func New(cfg *config.Config) *App {
	app := &App{Provider: NewProvider(cfg)}
	app.Server = server.New(
		server.WithLogger(app.Logger()),
		server.WithCors(app.Config().Cors.AllowOrigin),
		server.WithRender(render.NewFromFS(web.StaticFS, "dist/html")),
		server.WithValidate(validator.NewEchoValidator()),
		server.WithErrorHandle(handler.NewErrorHandler().Handle),
	)

	app.Server.ApplyController("", []server.Controller{
		app.WebAuth(),
		app.WebToken(),
	})

	app.Server.ApplyController("/api", []server.Controller{})

	app.Server.FileFS("/favicon.ico/", "dist/favicon.png", web.StaticFS)

	app.Closer().Add(app.Server.Shutdown)

	return app
}

func (app *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	addr := net.JoinHostPort(app.Config().Http.Host, app.Config().Http.Port)

	app.Logger().Info("server start", "addr", addr)

	go func(cancel context.CancelFunc) {
		if err := app.Server.Start(addr); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				app.Logger().Info("server stopped")
			} else {
				app.Logger().Error("server stopped", "error", err.Error())
			}
			cancel()
		}
	}(cancel)

	<-ctx.Done()

	return app.Close()
}

func (app *App) Close() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.Config().App.Shutdown)
	defer cancel()
	return app.Closer().Close(shutdownCtx)
}
