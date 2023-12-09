package config

type Config struct {
	App      App      `env:",prefix=APP_"`
	Logger   Logger   `env:",prefix=LOG_"`
	Http     Http     `env:",prefix=HTTP_"`
	Database Database `env:",prefix=DB_"`
	Path     Path     `env:",prefix=PATH_"`
}

func NewConfig() Config {
	return Config{}
}
