package config

import "net"

type Http struct {
	Host string `env:"HOST,default=0.0.0.0"`
	Port string `env:"PORT,default=80"`
}

func (c Http) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
