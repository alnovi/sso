package config

type Path struct {
	Web  string `env:"WEB,default=web"`
	Html string `env:"HTML,default=web/html"`
}
