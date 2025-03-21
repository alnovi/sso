package config

import "fmt"

type Database struct {
	Host     string `env:"HOST,default=localhost"`
	Port     string `env:"PORT,default=5432"`
	Username string `env:"USERNAME,default=root"`
	Password string `env:"PASSWORD,default=secret"`
	Database string `env:"DATABASE,default=sso"`
}

func (c *Database) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.Database)
}
