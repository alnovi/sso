package integration //nolint:typecheck

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/internal/config"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerLogger struct{}

func (l *ContainerLogger) Printf(_ string, _ ...interface{}) {}

type TestSuite struct {
	suite.Suite
	pgContainer   *postgres.PostgresContainer
	smtpContainer testcontainers.Container
	App           *server.App
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	ctx := context.Background()
	cfg := &config.Config{}
	s.initConfig(ctx, cfg)
	s.initDatabase(ctx, cfg)
	s.initMailServer(ctx, cfg)
	s.NotPanics(func() {
		s.App = server.New(cfg)
	})
}

func (s *TestSuite) TearDownSuite() {
	s.Require().NoError(s.App.Close())
	s.Require().NoError(s.pgContainer.Terminate(context.Background()))
}

func (s *TestSuite) SetupTest() {
	ctx := context.WithValue(context.Background(), config.KeyEnvironment, s.App.Config().App.Environment)
	ctx = context.WithValue(ctx, config.KeyClientAdminID, s.App.Config().Client.AdminID)
	ctx = context.WithValue(ctx, config.KeyClientProfileID, s.App.Config().Client.ProfileID)
	ctx = context.WithValue(ctx, config.KeyUserAdminID, s.App.Config().User.AdminID)
	ctx = context.WithValue(ctx, config.KeyUserAdminEmail, s.App.Config().User.AdminEmail)
	err := s.App.Repository().MigrateUp(ctx, s.App.Logger())
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownTest() {
	err := s.App.Repository().MigrateDown(context.Background(), s.App.Logger())
	s.Require().NoError(err)
}

func (s *TestSuite) BuildQuery(query map[string]string) string {
	q := make(url.Values)
	for k, v := range query {
		q.Set(k, v)
	}
	return q.Encode()
}

func (s *TestSuite) BuildFormData(mime string, data map[string]string) string {
	if mime == echo.MIMEApplicationForm {
		form := make(url.Values)
		for k, v := range data {
			form.Set(k, v)
		}
		return form.Encode()
	}

	form, _ := json.Marshal(data)
	return string(form)
}

func (s *TestSuite) SendToServer(h echo.HandlerFunc, c echo.Context) error {
	var err error
	if err = h(c); err != nil {
		s.App.Server.HTTPErrorHandler(err, c)
	}
	return err
}

func (s *TestSuite) initConfig(ctx context.Context, cfg *config.Config) {
	s.Require().NoError(configure.ParseEnv(ctx, cfg))
	cfg.App.Environment = config.EnvTesting
	cfg.Log.Format = logger.FormatStub
}

func (s *TestSuite) initDatabase(ctx context.Context, cfg *config.Config) {
	var err error
	s.pgContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithLogger(&ContainerLogger{}),
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase(cfg.DB.Database),
		postgres.WithUsername(cfg.DB.User),
		postgres.WithPassword(cfg.DB.Password),
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

	cfg.DB.Host = host
	cfg.DB.Port = port.Port()
}

func (s *TestSuite) initMailServer(ctx context.Context, cfg *config.Config) {
	var err error

	s.smtpContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Logger: &ContainerLogger{},
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mailhog/mailhog:latest",
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
