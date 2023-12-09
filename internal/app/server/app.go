package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/alnovi/sso/config"
	"github.com/labstack/echo/v4"
)

type App struct {
	cfg      *config.Config
	log      *slog.Logger
	adapters *Adapters
	http     *echo.Echo
}

func NewApp(cfg *config.Config, log *slog.Logger) (*App, error) {
	var err error

	app := &App{cfg: cfg, log: log}

	app.adapters, err = newAdapters(app)
	if err != nil {
		return nil, err
	}

	err = app.adapters.repo.MigrateUp()
	if err != nil {
		return nil, err
	}

	services, err := newServices(app, app.adapters)
	if err != nil {
		return nil, err
	}

	useCases, err := newUseCases(app, services)
	if err != nil {
		return nil, err
	}

	middlewares, err := newMiddlewares(app, services)
	if err != nil {
		return nil, err
	}

	handlers, err := newHandlers(app, useCases)
	if err != nil {
		return nil, err
	}

	app.http, err = NewHttpServer(app, middlewares, handlers)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	var err error

	go func() {
		addr := app.cfg.Http.Address()
		err = app.http.Start(addr)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				app.log.Info("server closed")
			} else {
				app.log.Error("failed to start server", "error", err.Error())
				os.Exit(1)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	err = app.http.Shutdown(context.Background())
	if err != nil {
		return err
	}

	err = app.adapters.repo.Close()
	if err != nil {
		return err
	}

	return nil
}
