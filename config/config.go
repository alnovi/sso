package config

import (
	"slices"

	"github.com/google/uuid"
)

const (
	AppEnvironmentProduction  = "production"
	AppEnvironmentDevelopment = "development"
	AppEnvironmentTesting     = "testing"
)

var Version = "v0.0.0"
var CtxConfigKey = "config"

type Config struct {
	App       App       `env:",prefix=APP_"`
	Logger    Logger    `env:",prefix=LOG_"`
	Http      Http      `env:",prefix=HTTP_"`
	Database  Database  `env:",prefix=DB_"`
	Mail      Mail      `env:",prefix=MAIL_"`
	Scheduler Scheduler `env:",prefix=SCHEDULER_"`
	Trace     Trace     `env:",prefix=TRACE_"`
	CAdmin    Client    `env:",prefix=CLIENT_ADMIN_"`
	UAdmin    User      `env:",prefix=USER_ADMIN_"`
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == AppEnvironmentProduction || !slices.Contains([]string{
		AppEnvironmentDevelopment,
		AppEnvironmentTesting,
	}, c.App.Environment)
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == AppEnvironmentDevelopment
}

func (c *Config) Normalize() {
	if err := uuid.Validate(c.UAdmin.Id); err != nil {
		c.UAdmin.Id = uuid.New().String()
	}
}
