package config

import "time"

type Scheduler struct {
	StopTimeout        time.Duration `env:"STOP_TIMEOUT,default=5s"`
	DeleteTokenExpired time.Duration `env:"DELETE_TOKEN_EXPIRED,default=5m"`
	DeleteSessionEmpty time.Duration `env:"DELETE_SESSION_EMPTY,default=5m"`
}
