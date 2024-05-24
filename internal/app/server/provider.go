package server

import (
	"context"
	"log/slog"

	"github.com/alnovi/sso/internal/adapter/repository/postgres"
	"github.com/alnovi/sso/internal/config"
	"github.com/alnovi/sso/internal/transport/http/handler/web"
	"github.com/alnovi/sso/internal/usecase/auth"
	"github.com/alnovi/sso/pkg/closer"
	"github.com/alnovi/sso/pkg/logger"
)

type Provider struct {
	cfg        *config.Config
	logger     *slog.Logger
	closer     *closer.Closer
	repository *postgres.Repository
	webAuth    *web.AuthHandler
}

func NewProvider(cfg *config.Config) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Config() *config.Config {
	return p.cfg
}

func (p *Provider) Logger() *slog.Logger {
	if p.logger == nil {
		p.logger = logger.New(p.cfg.Log.Format, p.cfg.Log.Level)
	}

	return p.logger
}

func (p *Provider) Closer() *closer.Closer {
	if p.closer == nil {
		p.closer = closer.New()
	}
	return p.closer
}

func (p *Provider) Repository() *postgres.Repository {
	var err error
	if p.repository == nil {
		cfg := p.cfg.DB

		if p.repository, err = postgres.New(cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.SSL); err != nil {
			panic(err)
		}

		ctx := context.WithValue(context.Background(), config.KeyEnvironment, p.Config().App.Environment)
		ctx = context.WithValue(ctx, config.KeyClientAdminID, p.Config().Client.AdminID)
		ctx = context.WithValue(ctx, config.KeyClientProfileID, p.Config().Client.ProfileID)
		ctx = context.WithValue(ctx, config.KeyUserAdminID, p.Config().User.AdminID)

		if err = p.repository.MigrateUp(ctx, p.Logger()); err != nil {
			panic(err)
		}

		p.Closer().Add(p.repository.Close)
	}

	return p.repository
}

func (p *Provider) WebAuth() *web.AuthHandler {
	if p.webAuth == nil {
		p.webAuth = web.NewAuthHandler(
			p.Config().Client.ProfileID,
			auth.NewUseCase(p.Repository()),
		)
	}
	return p.webAuth
}
