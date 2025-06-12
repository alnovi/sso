package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/alnovi/gomon/configure"
	"github.com/alnovi/gomon/logger"
	"github.com/alnovi/gomon/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/token"
)

const (
	TestIP             = "127.0.0.1"
	TestAgent          = "suite-test-agent"
	TestSecret         = "secret"
	TestRole           = entity.RoleManager
	ImagePostgres      = "postgres:16-alpine"
	ImageMailSMTP      = "mailhog/mailhog:latest"
	LoggerFormat       = logger.FormatDiscard
	LoggerLevel        = logger.LevelInfo
	MsgNotAssertCode   = "not assert code"
	MsgNotAssertBody   = "not assert body"
	MsgNotAssertHeader = "not assert header"
	MsgNotAssertError  = "not assert error"
)

var (
	TestUser   *entity.User
	TestClient *entity.Client
)

type ContainerLogger struct {
	logger *slog.Logger
}

func (l *ContainerLogger) Printf(f string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, args...))
}

func NewContainerLogger(format, level string) *ContainerLogger {
	return &ContainerLogger{logger: logger.New(logger.WithFormat(format), logger.WithLevel(level))}
}

type TestSuite struct {
	suite.Suite
	app           *server.App
	logger        *ContainerLogger
	pgContainer   *postgres.PostgresContainer
	smtpContainer testcontainers.Container
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	ctx := context.Background()
	cfg := &config.Config{}

	s.initConfig(ctx, cfg)
	s.initDockerLogger(ctx, cfg)
	s.initDatabase(ctx, cfg)
	s.initMailServer(ctx, cfg)
	s.initTestData(ctx, cfg)

	s.Require().NotPanics(func() {
		s.app = server.NewApp(cfg)
	})
}

func (s *TestSuite) TearDownSuite() {
	s.Require().NoError(s.app.Provider.Closer().Close())
	s.Require().NoError(s.pgContainer.Terminate(context.Background()))
	s.Require().NoError(s.smtpContainer.Terminate(context.Background()))
	_ = s.app.Provider.Certs().RemoveDir()
}

func (s *TestSuite) SetupTest() {
	ctx := context.Background()

	s.Require().NotPanics(func() {
		s.app.Provider.MigrationUp()
	})

	err := s.app.Provider.Repository().ClientCreate(ctx, TestClient)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().UserCreate(ctx, TestUser)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().RoleUpdate(ctx, &entity.Role{ClientId: TestClient.Id, UserId: TestUser.Id, Role: TestRole})
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownTest() {
	s.Require().NotPanics(func() {
		s.app.Provider.MigrationDown()
	})
}

func (s *TestSuite) initConfig(ctx context.Context, cfg *config.Config) {
	s.Require().NoError(configure.LoadFromEnv(ctx, cfg))

	cfg.App.Environment = config.AppEnvironmentTesting
	cfg.App.Host = "http://localhost:8080"

	cfg.Logger.Format = LoggerFormat
	cfg.Logger.Level = LoggerLevel

	cfg.Database.Username = "postgres"
	cfg.Database.Password = "postgres"
	cfg.Database.Database = "sso"

	cfg.Mail.Username = "test@example.com"
	cfg.Mail.Password = TestSecret

	cfg.CAdmin.Id = "sso-admin"
	cfg.CAdmin.Name = "Client Admin"
	cfg.CAdmin.Secret = TestSecret
	cfg.CAdmin.Callback = "/admin/callback"

	cfg.UAdmin.Name = "User Admin"
	cfg.UAdmin.Email = "admin@example.com"
	cfg.UAdmin.Password = TestSecret

	cfg.Normalize()
}

func (s *TestSuite) initDockerLogger(_ context.Context, cfg *config.Config) {
	s.logger = NewContainerLogger(cfg.Logger.Format, cfg.Logger.Level)
}

func (s *TestSuite) initDatabase(ctx context.Context, cfg *config.Config) {
	var err error
	s.pgContainer, err = postgres.Run(ctx, ImagePostgres,
		testcontainers.WithLogger(s.logger),
		postgres.WithDatabase(cfg.Database.Database),
		postgres.WithUsername(cfg.Database.Username),
		postgres.WithPassword(cfg.Database.Password),
		testcontainers.WithWaitStrategy(wait.
			ForLog("database system is ready").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second),
		),
	)
	s.Require().NoError(err)

	host, err := s.pgContainer.Host(ctx)
	s.Require().NoError(err)

	port, err := s.pgContainer.MappedPort(ctx, "5432/tcp")
	s.Require().NoError(err)

	cfg.Database.Host = host
	cfg.Database.Port = port.Port()
}

func (s *TestSuite) initMailServer(ctx context.Context, cfg *config.Config) {
	var err error

	s.smtpContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Logger: s.logger,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        ImageMailSMTP,
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		},
		Started: true,
	})
	s.Require().NoError(err)

	host, err := s.smtpContainer.Host(ctx)
	s.Require().NoError(err)

	port, err := s.smtpContainer.MappedPort(ctx, "1025/tcp")
	s.Require().NoError(err)

	cfg.Mail.Host = host
	cfg.Mail.Port = port.Port()
}

func (s *TestSuite) initTestData(_ context.Context, _ *config.Config) {
	password, _ := utils.HashPassword("password")

	TestUser = &entity.User{
		Id:       uuid.NewString(),
		Name:     "Test user",
		Email:    "test@example.com",
		Password: password,
	}

	TestClient = &entity.Client{
		Id:       "test-client",
		Name:     "Test client",
		Secret:   TestSecret,
		Callback: "http://localhost/callback",
		IsSystem: false,
	}
}

func (s *TestSuite) config() *config.Config {
	return s.app.Provider.Config()
}

func (s *TestSuite) buildQuery(query map[string]string) string {
	q := make(url.Values)
	for k, v := range query {
		q.Set(k, v)
	}
	return q.Encode()
}

func (s *TestSuite) applyHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (s *TestSuite) applyCookies(req *http.Request, cookies []*http.Cookie) {
	for _, c := range cookies {
		req.AddCookie(c)
	}
}

func (s *TestSuite) buildData(mime string, data map[string]any) string {
	switch mime {
	case echo.MIMEApplicationJSON:
		return s.buildDataJson(data)
	default:
		return s.buildDataForm(data)
	}
}

func (s *TestSuite) buildDataJson(data map[string]any) string {
	body, _ := json.Marshal(data)
	return string(body)
}

func (s *TestSuite) buildDataForm(data map[string]any) string {
	form := make(url.Values)
	for k, v := range data {
		form.Set(k, fmt.Sprintf("%v", v))
	}
	return form.Encode()
}

func (s *TestSuite) sendToServer(h echo.HandlerFunc, c echo.Context, mws ...echo.MiddlewareFunc) error {
	var err error
	var mwh echo.HandlerFunc

	for _, mw := range mws {
		mwh = mw(func(c echo.Context) error {
			return nil
		})

		if err = mwh(c); err != nil {
			s.app.HttpServer.HTTPErrorHandler(err, c)
			return err
		}
	}

	if err = h(c); err != nil {
		s.app.HttpServer.HTTPErrorHandler(err, c)
	}

	return err
}

func (s *TestSuite) sendToMiddleware(c echo.Context, ms ...echo.MiddlewareFunc) error {
	if len(ms) == 0 {
		return errors.New("middlewares must not be empty")
	}

	for _, m := range ms {
		h := m(func(c echo.Context) error {
			return nil
		})

		if err := h(c); err != nil {
			s.app.HttpServer.HTTPErrorHandler(err, c)
			return err
		}
	}

	return nil
}

func (s *TestSuite) accessTokens(clientId, userId, role string, opts ...token.Option) (session *entity.Session, access, refresh *entity.Token, err error) {
	session = &entity.Session{
		Id:     uuid.NewString(),
		UserId: userId,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	if err = s.app.Provider.Repository().SessionCreate(context.Background(), session); err != nil {
		return
	}

	if access, err = s.app.Provider.Token().AccessToken(context.Background(), session.Id, clientId, userId, "User name", role, opts...); err != nil {
		return
	}

	if refresh, err = s.app.Provider.Token().RefreshToken(context.Background(), session.Id, clientId, userId, time.Now(), opts...); err != nil {
		return
	}

	return
}
