package config

import (
	"time"
)

const (
	AppEnvironmentProduction  = "production"
	AppEnvironmentDevelopment = "development"
	AppEnvironmentTesting     = "testing"
)

type Config struct {
	App        App      `env:",prefix=APP_"`
	Logger     Logger   `env:",prefix=LOG_"`
	Mail       Mail     `env:",prefix=MAIL_"`
	Http       Http     `env:",prefix=HTTP_"`
	Cookie     Cookie   `env:",prefix=COOKIE_"`
	Jwt        Jwt      `env:",prefix=JWT_"`
	Database   Database `env:",prefix=DB_"`
	Client     Client   `env:",prefix=CLIENT_"`
	Admin      User     `env:",prefix=ADMIN_"`
	TestClient Client
	TestUser   User
}

type App struct {
	Environment string        `env:"ENVIRONMENT,default=production"`
	Version     string        `env:"VERSION,default=0.0.1"`
	Host        string        `env:"HOST,required"`
	Secret      string        `env:"SECRET,required"`
	Shutdown    time.Duration `env:"SHUTDOWN,default=5s"`
}

type Logger struct {
	Format string `env:"FORMAT,default=json"`
	Level  string `env:"LEVEL,default=error"`
}

type Http struct {
	Host string `env:"HOST,default=0.0.0.0"`
	Port string `env:"PORT,default=8080"`
}

type Cookie struct {
	Secure bool `env:"SECURE,default=true"`
}

type Database struct {
	Host     string `env:"HOST,default=localhost"`
	Port     string `env:"PORT,default=5432"`
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
	Database string `env:"DATABASE,required"`
}

type Mail struct {
	Host     string `env:"HOST,default=smtp.gmail.com"`
	Port     string `env:"PORT,default=587"`
	From     string `env:"FROM,default=SSO"`
	Username string `env:"USERNAME,required"`
	Password string `env:"PASSWORD,required"`
}

type Jwt struct {
	PrivateKey string `env:"PRIVATE_KEY,required"`
	PublicKey  string `env:"PUBLIC_KEY,required"`
}

type Client struct {
	Id     string `env:"ID,default=admin"`
	Name   string `env:"NAME,default=Админ панель"`
	Secret string `env:"SECRET,default=secret"`
	Host   string `env:"HOST,default=https://127.0.0.1"`
}

type User struct {
	Id       string `env:"ID,default=00000000-0000-0000-0000-000000000000"`
	Name     string `env:"NAME,default=Admin"`
	Email    string `env:"EMAIL,default=admin@example.com"`
	Password string `env:"PASSWORD,default=secret"`
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == AppEnvironmentProduction
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == AppEnvironmentDevelopment
}

func (c *Config) IsTesting() bool {
	return c.App.Environment == AppEnvironmentTesting
}
