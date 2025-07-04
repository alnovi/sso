package provider

import (
	"context"
	"log/slog"

	"github.com/alnovi/gomon/closer"
	"github.com/alnovi/gomon/configure"
	"github.com/alnovi/gomon/logger"
	"github.com/alnovi/gomon/migrator"
	"github.com/alnovi/gomon/utils"
	"github.com/alnovi/gomon/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.30.0"
	itrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/adapter/mailing"
	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/helper"
	"github.com/alnovi/sso/internal/service/admin"
	"github.com/alnovi/sso/internal/service/certs"
	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/crontask"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/service/profile"
	"github.com/alnovi/sso/internal/service/rule"
	"github.com/alnovi/sso/internal/service/stats"
	"github.com/alnovi/sso/internal/service/storage"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/pkg/database/postgres"
	"github.com/alnovi/sso/pkg/scheduler"
	_ "github.com/alnovi/sso/scripts/migrations"
)

type Provider struct {
	config      *config.Config
	logger      *slog.Logger
	tracer      itrace.TracerProvider
	closer      *closer.Closer
	validator   *validator.Validator
	db          *postgres.Client
	migrator    *migrator.Migrator
	repository  *repository.Repository
	transaction repository.Transaction
	mailing     *mailing.Mailing
	scheduler   *scheduler.Scheduler
	certs       *certs.Certs
	token       *token.Token
	oauth       *oauth.OAuth
	cookie      *cookie.Cookie
	profile     *profile.UserProfile
	admin       *admin.Admin
	clients     *storage.Clients
	users       *storage.Users
	roles       *storage.Roles
	sessions    *storage.Sessions
	stats       *stats.Stats
}

func New(config *config.Config) *Provider {
	return &Provider{config: config}
}

func (p *Provider) Config() *config.Config {
	if p.config == nil {
		p.config = new(config.Config)

		err := configure.LoadFromEnv(context.Background(), p.config)
		utils.MustMsg(err, "failed to load environment variables config")

		p.config.Normalize()
	}
	return p.config
}

func (p *Provider) Logger() *slog.Logger {
	if p.logger == nil {
		p.logger = logger.New(
			logger.WithFormat(p.Config().Logger.Format),
			logger.WithLevel(p.Config().Logger.Level),
		)
	}
	return p.logger
}

func (p *Provider) LoggerMod(mod string) *slog.Logger {
	if mod == "" {
		return p.Logger()
	}
	return p.Logger().With("module", mod)
}

func (p *Provider) Closer() *closer.Closer {
	if p.closer == nil {
		p.closer = closer.New(p.Config().App.Shutdown)
	}
	return p.closer
}

func (p *Provider) Tracer() itrace.TracerProvider {
	if p.tracer == nil {
		ctx := context.Background()

		if !p.Config().Trace.Enable || p.Config().Trace.ExportAddr == "" {
			p.tracer = noop.NewTracerProvider()
			otel.SetTracerProvider(p.tracer)
			return p.tracer
		}

		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		exporter, err := helper.NewGrpcExporter(ctx, p.Config().Trace.ExportAddr)
		utils.MustMsg(err, "failed to create tracer exporter")

		resourceOptions := []resource.Option{
			resource.WithAttributes(semconv.ServiceName(helper.TraceServiceName)),
			resource.WithHost(),
			resource.WithTelemetrySDK(),
		}

		resources, err := resource.New(ctx, resourceOptions...)
		utils.MustMsg(err, "failed to create tracer resources")

		tp := trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithBatcher(exporter, trace.WithBatchTimeout(p.Config().Trace.BatchTimeout)),
			trace.WithResource(resources),
		)

		p.Closer().Add(tp.Shutdown)

		p.tracer = tp

		otel.SetTracerProvider(p.tracer)
	}
	return p.tracer
}

func (p *Provider) Validator() *validator.Validator {
	if p.validator == nil {
		var err error

		p.validator = validator.NewValidator()

		err = p.validator.AddRule(rule.NewClientID())
		utils.MustMsg(err, "failed to add rule 'client id'")
	}
	return p.validator
}

func (p *Provider) DB() *postgres.Client {
	if p.db == nil {
		var err error

		p.db, err = postgres.NewClient(p.Config().Database.DSN(), postgres.WithLogger(p.LoggerMod("sql")))
		utils.MustMsg(err, "failed to connect to database")

		err = p.db.Ping(context.Background())
		utils.MustMsg(err, "failed to ping database")

		p.Closer().Add(func(_ context.Context) error {
			return p.db.Close()
		})
	}
	return p.db
}

func (p *Provider) Repository() *repository.Repository {
	if p.repository == nil {
		p.repository = repository.NewRepository(p.DB())
	}
	return p.repository
}

func (p *Provider) Transaction() repository.Transaction {
	if p.transaction == nil {
		p.transaction = postgres.NewTransaction(p.DB().Master())
	}
	return p.transaction
}

func (p *Provider) Migrator() *migrator.Migrator {
	if p.migrator == nil {
		p.migrator = migrator.NewMigrator(
			migrator.WithLogger(p.LoggerMod("migrator")),
			migrator.WithDialect(migrator.DialectPostgres),
		)
	}
	return p.migrator
}

func (p *Provider) MigrationUp() {
	ctx := context.WithValue(context.Background(), config.CtxConfigKey, p.Config())

	db := p.DB().DB()
	defer func() {
		_ = db.Close()
	}()

	err := p.Migrator().UpContext(ctx, db)
	utils.Must(err)
}

func (p *Provider) MigrationDown() {
	ctx := context.WithValue(context.Background(), config.CtxConfigKey, p.Config())

	db := p.DB().DB()
	defer func() {
		_ = db.Close()
	}()

	err := p.Migrator().ResetContext(ctx, db)
	utils.Must(err)
}

func (p *Provider) Mailing() *mailing.Mailing {
	if p.mailing == nil {
		var err error

		p.mailing, err = mailing.New(
			p.Config().Mail.Host,
			p.Config().Mail.Port,
			mailing.WithAppHost(p.Config().App.Host),
			mailing.WithFrom(p.Config().Mail.From, p.Config().Mail.Username),
			mailing.WithAuthUsername(p.Config().Mail.Username),
			mailing.WithAuthPassword(p.Config().Mail.Password),
		)

		utils.MustMsg(err, "failed to connect to mailing service")

		utils.MustMsg(p.mailing.Ping(context.Background()), "failed to ping mailing service")

		p.Closer().Add(p.mailing.Close)
	}
	return p.mailing
}

func (p *Provider) Scheduler() *scheduler.Scheduler {
	if p.scheduler == nil {
		var err error

		p.scheduler, err = scheduler.New(p.Config().Scheduler.StopTimeout)
		utils.MustMsg(err, "failed create scheduler")

		err = p.scheduler.AddDurationTask(p.Config().Scheduler.DeleteTokenExpired, crontask.NewTaskDeleteTokenExpired(p.Repository()))
		utils.MustMsg(err, "failed add delete token expired task")

		err = p.scheduler.AddDurationTask(p.Config().Scheduler.DeleteSessionEmpty, crontask.NewTaskDeleteSessionEmpty(p.Repository()))
		utils.MustMsg(err, "failed add delete session empty task")

		p.Closer().Add(func(_ context.Context) error {
			return p.scheduler.Stop()
		})
	}
	return p.scheduler
}

func (p *Provider) Certs() *certs.Certs {
	if p.certs == nil {
		var err error
		p.certs, err = certs.New()
		utils.MustMsg(err, "failed init certs service")
	}
	return p.certs
}

func (p *Provider) Token() *token.Token {
	if p.token == nil {
		publicKey, privateKey, err := p.Certs().Keys()
		utils.MustMsg(err, "failed get rsa keys")

		p.token, err = token.New(privateKey, publicKey, p.Repository())
		utils.MustMsg(err, "failed to init Token service")
	}
	return p.token
}

func (p *Provider) OAuth() *oauth.OAuth {
	if p.oauth == nil {
		p.oauth = oauth.NewOAuth(p.Repository(), p.Transaction(), p.Token(), p.Mailing())
	}
	return p.oauth
}

func (p *Provider) Cookie() *cookie.Cookie {
	if p.cookie == nil {
		p.cookie = cookie.New(p.Config().IsProduction())
	}
	return p.cookie
}

func (p *Provider) Profile() *profile.UserProfile {
	if p.profile == nil {
		p.profile = profile.NewUserProfile(p.Repository(), p.Transaction())
	}
	return p.profile
}

func (p *Provider) Admin() *admin.Admin {
	if p.admin == nil {
		p.admin = admin.NewAdmin(p.Config().CAdmin.Id, p.Repository(), p.Transaction(), p.OAuth())
	}
	return p.admin
}

func (p *Provider) StorageClients() *storage.Clients {
	if p.clients == nil {
		p.clients = storage.NewClients(p.Repository(), p.Transaction())
	}
	return p.clients
}

func (p *Provider) StorageUsers() *storage.Users {
	if p.users == nil {
		p.users = storage.NewUsers(p.Repository(), p.Transaction())
	}
	return p.users
}

func (p *Provider) StorageRoles() *storage.Roles {
	if p.roles == nil {
		p.roles = storage.NewRoles(p.Repository(), p.Transaction())
	}
	return p.roles
}

func (p *Provider) StorageSessions() *storage.Sessions {
	if p.sessions == nil {
		p.sessions = storage.NewSessions(p.Repository(), p.Transaction())
	}
	return p.sessions
}

func (p *Provider) Stats() *stats.Stats {
	if p.stats == nil {
		p.stats = stats.NewStats(p.Repository())
	}
	return p.stats
}
