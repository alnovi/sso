package config

type User struct {
	Id       string `env:"ID,default=00000000-0000-0000-0000-000000000000"`
	Name     string `env:"NAME,default=Admin"`
	Email    string `env:"EMAIL,default=admin@example.com"`
	Password string `env:"PASSWORD,default=secret"`
}
