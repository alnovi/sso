package config

import "time"

type App struct {
	Environment string        `env:"ENVIRONMENT,default=production"`
	Version     string        `env:"VERSION"`
	Host        string        `env:"HOST"`
	Secret      string        `env:"SECRET,default=XbZUD6rE49aMPpEB"`
	Shutdown    time.Duration `env:"SHUTDOWN,default=5s"`
}
