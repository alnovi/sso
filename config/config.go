package config

const (
	AppEnvironmentProduction  = "production"
	AppEnvironmentDevelopment = "development"
	AppEnvironmentTesting     = "testing"
)

type Config struct {
	App       App       `env:",prefix=APP_"`
	Logger    Logger    `env:",prefix=LOG_"`
	Http      Http      `env:",prefix=HTTP_"`
	Database  Database  `env:",prefix=DB_"`
	Mail      Mail      `env:",prefix=MAIL_"`
	Scheduler Scheduler `env:",prefix=SCHEDULER_"`
	Jwt       Jwt       `env:",prefix=JWT_"`
	CAdmin    Client    `env:",prefix=CLIENT_ADMIN_"`
	UAdmin    User      `env:",prefix=USER_ADMIN_"`
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
