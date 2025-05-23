package config

import "github.com/google/uuid"

type User struct {
	Id       string `env:"ID"`
	Name     string `env:"NAME,default=Admin"`
	Email    string `env:"EMAIL,default=admin@example.com"`
	Password string `env:"PASSWORD,default=secret"`
}

func (c User) ID() string {
	if err := uuid.Validate(c.Id); err != nil {
		c.Id = uuid.New().String()
	}
	return c.Id
}
