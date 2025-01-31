package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alnovi/sso/config"
	pgrepo "github.com/alnovi/sso/internal/adapter/repository/postgres"
	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/jwt"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/service/sessions"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/pkg/closer"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/db/pgs"
	"github.com/alnovi/sso/pkg/logger"
	"github.com/alnovi/sso/pkg/migrator"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
	_ "github.com/alnovi/sso/scripts/migrations"
)

type Provider struct {
	config      *config.Config
	logger      *slog.Logger
	closer      *closer.Closer
	validator   *validator.EchoValidator
	dbPool      *pgs.Client
	transaction *pgs.Transaction
	repository  *pgrepo.Repository
	cookie      *cookie.Cookie
	oauth       *oauth.OAuth
	jwt         *jwt.JWT
	session     *sessions.Session
	token       *token.Token
}

func New(cfg *config.Config) *Provider {
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

func (p *Provider) DB() *pgs.Client {
	if p.dbPool == nil {
		var err error

		cfg := p.Config().Database
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

		p.dbPool, err = pgs.NewClient(dsn)
		if err != nil {
			utils.Must(fmt.Errorf("failed to connect to database: %w", err))
		}

		if err = p.dbPool.Ping(context.Background()); err != nil {
			utils.Must(fmt.Errorf("failed to ping database: %w", err))
		}

		p.dbPool.SetLogger(p.LoggerModule("sql"))

		p.Closer().Add(func(_ context.Context) error {
			return p.dbPool.Close()
		})
	}
	return p.dbPool
}

func (p *Provider) Transaction() *pgs.Transaction {
	if p.transaction == nil {
		p.transaction = pgs.NewTransaction(p.DB().DB())
	}
	return p.transaction
}

func (p *Provider) Repository() *pgrepo.Repository {
	if p.repository == nil {
		p.repository = pgrepo.New(p.DB())
	}
	return p.repository
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

func (p *Provider) Cookie() *cookie.Cookie {
	if p.cookie == nil {
		p.cookie = cookie.New(p.Config().IsProduction())
	}
	return p.cookie
}

func (p *Provider) OAuth() *oauth.OAuth {
	if p.oauth == nil {
		p.oauth = oauth.New(p.Repository(), p.Transaction(), p.Token(), p.Session())
	}
	return p.oauth
}

func (p *Provider) JWT() *jwt.JWT {
	if p.jwt == nil {
		var err error
		p.jwt, err = jwt.New([]byte(p.Config().Jwt.PrivateKey), []byte(p.Config().Jwt.PublicKey))
		utils.Must(err)
	}
	return p.jwt
}

func (p *Provider) Session() *sessions.Session {
	if p.session == nil {
		p.session = sessions.New(p.Repository())
	}
	return p.session
}

func (p *Provider) Token() *token.Token {
	if p.token == nil {
		p.token = token.New(p.Repository(), p.JWT())
	}
	return p.token
}
