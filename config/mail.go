package config

type Mail struct {
	Host     string `env:"HOST,default=smtp.gmail.com"`
	Port     string `env:"PORT,default=587"`
	From     string `env:"FROM,default=SSO"`
	Username string `env:"USERNAME,default=sso@example.com"`
	Password string `env:"PASSWORD,default=secret"`
}
