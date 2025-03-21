package config

type Client struct {
	Id       string `env:"ID,default=sso-admin"`
	Name     string `env:"NAME,default=Пользователи"`
	Secret   string `env:"SECRET,default=secret"`
	Callback string `env:"CALLBACK,default=http://127.0.0.1/admin/callback"`
}
