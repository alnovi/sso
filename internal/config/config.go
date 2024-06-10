package config

import "time"

const (
	EnvProduction = "production"
	EnvTesting    = "testing"

	KeyEnvironment key = iota
	KeyClientAdminID
	KeyClientProfileID
	KeyUserAdminID
	KeyUserAdminEmail
)

type key int

type Config struct {
	App    App    `env:",prefix=APP_"`
	Log    Log    `env:",prefix=LOG_"`
	Mail   Mail   `env:",prefix=MAIL_"`
	DB     DB     `env:",prefix=DB_"`
	Http   Http   `env:",prefix=HTTP_"`
	Cors   Cors   `env:",prefix=CORS_"`
	Client Client `env:",prefix=CLIENT_"`
	User   User   `env:",prefix=USER_"`
}

type App struct {
	Environment string        `env:"ENVIRONMENT,default=production"`
	Host        string        `env:"HOST"`
	Shutdown    time.Duration `env:"SHUTDOWN,default=5s"`
}

type Log struct {
	Format string `env:"FORMAT,default=json"`
	Level  string `env:"LEVEL,default=error"`
}

type Mail struct {
	FromName string `env:"NAME,default=SSO"`
	FromAddr string `env:"FROM,default=sso@example.com"`
	Host     string `env:"HOST,default=localhost"`
	Port     string `env:"PORT,default=5432"`
	User     string `env:"USER,default=example"`
	Password string `env:"PASSWORD,default=example"`
}

type DB struct {
	Host     string `env:"HOST,default=localhost"`
	Port     string `env:"PORT,default=5432"`
	Database string `env:"DATABASE,default=example"`
	User     string `env:"USER,default=example"`
	Password string `env:"PASSWORD,default=secret"`
	SSL      bool   `env:"SSL,default=false"`
}

type Http struct {
	Host string `env:"HOST,default=0.0.0.0"`
	Port string `env:"PORT,default=8080"`
}

type Cors struct {
	AllowOrigin string `env:"ALLOW_ORIGIN"`
}

type Client struct {
	AdminID   string `env:"ADMIN_ID,default=00000000-0000-0000-0000-000000000001"`
	ProfileID string `env:"PROFILE_ID,default=00000000-0000-0000-0000-000000000002"`
}

type User struct {
	AdminID    string `env:"ADMIN_ID,default=00000000-0000-0000-0000-000000000001"`
	AdminEmail string `env:"ADMIN_EMAIL,default=admin@example.com"`
}
