package config

import "time"

type Trace struct {
	Enable       bool          `env:"ENABLE,default=false"`
	ExportAddr   string        `env:"EXPORT_ADDR,default=127.0.0.1:4317"`
	BatchTimeout time.Duration `env:"BATCH_TIMEOUT,default=5s"`
}
