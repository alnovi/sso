package config

type User struct {
	Id       string `env:"ID"`
	Name     string `env:"NAME,default=Admin"`
	Email    string `env:"EMAIL,default=admin@example.com"`
	Password string `env:"PASSWORD,default=secret"`
}
