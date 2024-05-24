package integration //nolint:typecheck

import (
	"context"
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
	pgContainer *postgres.PostgresContainer
	App         *server.App
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	ctx := context.Background()
	cfg := &config.Config{}
	s.initConfig(ctx, cfg)
	s.initDatabase(ctx, cfg)
	s.NotPanics(func() {
		s.App = server.New(cfg)
	})
}

func (s *TestSuite) TearDownSuite() {
	s.Must(s.App.Close())
	s.Must(s.pgContainer.Terminate(context.Background()))
}

func (s *TestSuite) SetupTest() {
	ctx := context.WithValue(context.Background(), config.KeyEnvironment, s.App.Config().App.Environment)
	ctx = context.WithValue(ctx, config.KeyClientAdminID, s.App.Config().Client.AdminID)
	ctx = context.WithValue(ctx, config.KeyClientProfileID, s.App.Config().Client.ProfileID)
	ctx = context.WithValue(ctx, config.KeyUserAdminID, s.App.Config().User.AdminID)
	err := s.App.Repository().MigrateUp(ctx, s.App.Logger())
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownTest() {
	err := s.App.Repository().MigrateDown(context.Background(), s.App.Logger())
	s.Require().NoError(err)
}

func (s *TestSuite) Must(err error) {
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *TestSuite) SendToServer(h echo.HandlerFunc, c echo.Context) error {
	var err error
	if err = h(c); err != nil {
		s.App.Server.HTTPErrorHandler(err, c)
	}
	return err
}

func (s *TestSuite) initConfig(ctx context.Context, cfg *config.Config) {
	s.Must(configure.ParseEnv(ctx, cfg))
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
	s.Must(err)

	host, err := s.pgContainer.Host(ctx)
	s.Require().NoError(err)

	port, err := s.pgContainer.MappedPort(ctx, "5432/tcp")
	s.Require().NoError(err)

	cfg.DB.Host = host
	cfg.DB.Port = port.Port()
}
