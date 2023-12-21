package config

type Path struct {
	Html   string `env:"HTML,default=web/html"`
	Assets string `env:"ASSETS,default=web/assets"`
	Store  string `env:"STORE,default=web/public"`
}
