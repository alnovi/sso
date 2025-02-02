package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/logger"
)

const (
	TestIP            = "127.0.0.1"
	TestAgent         = "suite-test-agent"
	TestSecret        = "secret"
	TestRoleAdmin     = "admin"
	ImagePostgres     = "postgres:15.3-alpine"
	ImageMailSMTP     = "mailhog/mailhog:latest"
	LoggerFormat      = logger.FormatJson
	LoggerLevel       = logger.LevelError
	MsgNotAssertCode  = "not assert code"
	MsgNotAssertBody  = "not assert body"
	MsgNotAssertError = "not assert error"
)

type ContainerLogger struct{}

func (l *ContainerLogger) Printf(f string, args ...interface{}) {
	//fmt.Printf(f+"\n", args...)
}

type TestSuite struct {
	suite.Suite
	app           *server.App
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
	s.initRsaKeys(ctx, cfg)
	s.initDatabase(ctx, cfg)
	s.initMailServer(ctx, cfg)

	s.Require().NotPanics(func() {
		s.app = server.NewApp(cfg)
	})
}

func (s *TestSuite) TearDownSuite() {
	s.Require().NoError(s.app.Provider.Closer().Close())
	s.Require().NoError(s.pgContainer.Terminate(context.Background()))
}

func (s *TestSuite) SetupTest() {
	s.Require().NotPanics(func() {
		s.app.Provider.MigrationUp(context.Background())
	})
}

func (s *TestSuite) TearDownTest() {
	s.Require().NotPanics(func() {
		s.app.Provider.MigrationDown(context.Background())
	})
}

func (s *TestSuite) initConfig(_ context.Context, cfg *config.Config) {
	cfg.App.Environment = config.AppEnvironmentTesting
	cfg.App.Host = "http://localhost:8080"
	cfg.App.Secret = "1234567890AbCdFg"

	cfg.Logger.Format = LoggerFormat
	cfg.Logger.Level = LoggerLevel

	cfg.Database.Username = "postgres"
	cfg.Database.Password = "postgres"
	cfg.Database.Database = "sso"

	cfg.Mail.Username = "test@example.com"
	cfg.Mail.Password = TestSecret

	cfg.Jwt.PrivateKey = TestSecret
	cfg.Jwt.PublicKey = TestSecret

	cfg.TestClient.Id = uuid.NewString()
	cfg.TestClient.Name = "Test client"
	cfg.TestClient.Secret = TestSecret
	cfg.TestClient.Host = "https://127.0.0.1"

	cfg.TestUser.Id = uuid.NewString()
	cfg.TestUser.Name = "Test user"
	cfg.TestUser.Email = "test@example.com"
	cfg.TestUser.Password = TestSecret

	s.Require().NoError(configure.ParseEnv(cfg))
}

func (s *TestSuite) initRsaKeys(_ context.Context, cfg *config.Config) {
	buf := bytes.NewBuffer(nil)

	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Errorf("Cannot generate RSA key: %s\n", err))
	}

	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}

	err = pem.Encode(buf, privateKey)
	if err != nil {
		panic(fmt.Errorf("Cannot encode RSA key: %s\n", err))
	}

	cfg.Jwt.PrivateKey = buf.String()

	buf = bytes.NewBuffer(nil)

	public, err := asn1.Marshal(private.PublicKey)
	if err != nil {
		panic(fmt.Errorf("Cannot marshal RSA key: %s\n", err))
	}

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: public,
	}

	err = pem.Encode(buf, pemKey)
	if err != nil {
		panic(fmt.Errorf("Cannot encode public key: %s\n", err))
	}

	cfg.Jwt.PublicKey = buf.String()
}

func (s *TestSuite) initDatabase(ctx context.Context, cfg *config.Config) {
	var err error
	s.pgContainer, err = postgres.Run(ctx, ImagePostgres,
		testcontainers.WithLogger(&ContainerLogger{}),
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
		Logger: &ContainerLogger{},
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

func (s *TestSuite) buildQuery(query map[string]string) string {
	q := make(url.Values)
	for k, v := range query {
		q.Set(k, v)
	}
	return q.Encode()
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
			return c.String(http.StatusOK, "test")
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
