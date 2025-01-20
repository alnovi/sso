package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/pkg/client/postgres"
	"github.com/alnovi/sso/pkg/closer"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/logger"
	"github.com/alnovi/sso/pkg/migrator"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
	_ "github.com/alnovi/sso/scripts/migrations"
)

type Provider struct {
	config    *config.Config
	logger    *slog.Logger
	closer    *closer.Closer
	validator *validator.EchoValidator
	db        *postgres.Client
	tm        *postgres.Transaction
}

func NewProvider(cfg *config.Config) *Provider {
	return &Provider{config: cfg}
}

func (p *Provider) Config() *config.Config {
	if p.config == nil {
		p.config = new(config.Config)
		err := configure.ParseEnv(p.config)
		utils.Must(err)
	}
	return p.config
}

func (p *Provider) Logger() *slog.Logger {
	if p.logger == nil {
		p.logger = logger.New(p.Config().Logger.Format, p.Config().Logger.Level)
	}
	return p.logger
}

func (p *Provider) LoggerModule(module string) *slog.Logger {
	return p.Logger().With("module", module)
}

func (p *Provider) Closer() *closer.Closer {
	if p.closer == nil {
		p.closer = closer.New(p.Config().App.Shutdown)
	}
	return p.closer
}

func (p *Provider) Validator() *validator.EchoValidator {
	if p.validator == nil {
		p.validator = validator.NewEchoValidator()
	}
	return p.validator
}

func (p *Provider) DB() *postgres.Client {
	if p.db == nil {
		var err error

		cfg := p.Config().Database
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

		p.db, err = postgres.NewClient(dsn)
		if err != nil {
			utils.Must(fmt.Errorf("failed to connect to database: %w", err))
		}

		if err = p.db.Ping(context.Background()); err != nil {
			utils.Must(fmt.Errorf("failed to ping database: %w", err))
		}

		p.db.SetLogger(p.LoggerModule("sql"))

		p.Closer().Add(func(_ context.Context) error {
			return p.db.Close()
		})
	}
	return p.db
}

func (p *Provider) Transaction() *postgres.Transaction {
	if p.tm == nil {
		p.tm = postgres.NewTransaction(p.DB().DB())
	}
	return p.tm
}

func (p *Provider) MigrationUp(ctx context.Context) {
	ctx = context.WithValue(ctx, migrator.ConfigKey, p.Config())
	log := migrator.NewGooseLogger(p.LoggerModule("migrate"))
	db := p.DB().SqlDB()

	defer func() {
		_ = db.Close()
	}()

	err := migrator.PostgresUpFromPath(ctx, db, ".", log)
	utils.Must(err)
}

func (p *Provider) MigrationDown(ctx context.Context) {
	ctx = context.WithValue(ctx, migrator.ConfigKey, p.Config())
	log := migrator.NewGooseLogger(p.LoggerModule("migrate"))
	db := p.DB().SqlDB()

	defer func() {
		_ = db.Close()
	}()

	err := migrator.PostgresDownFromPath(ctx, db, ".", log)
	utils.Must(err)
}
