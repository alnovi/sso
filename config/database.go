package config

type Database struct {
	Host     string `env:"HOST,default=localhost"`
	Port     string `env:"PORT,default=5432"`
	Database string `env:"DATABASE,default=example"`
	User     string `env:"USER,default=example"`
	Password string `env:"PASSWORD,default=example"`
	SSLMode  bool   `env:"SSL,default=true"`
}

func (c Database) DSN() string {
	return ""
}
