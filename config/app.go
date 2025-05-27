package config

import "time"

type App struct {
	Environment string        `env:"ENVIRONMENT,default=production"`
	Host        string        `env:"HOST"`
	Shutdown    time.Duration `env:"SHUTDOWN,default=10s"`
}
