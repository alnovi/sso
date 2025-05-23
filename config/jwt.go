package config

import "path/filepath"

type Jwt struct {
	CertsDir   string `env:"CERTS_DIR,default=./certs"`
	PrivateKey string `env:"PRIVATE_KEY,default=private.pem"`
	PublicKey  string `env:"PUBLIC_KEY,default=public.pem"`
}

func (c Jwt) PrivatePath() string {
	return filepath.Join(c.CertsDir, c.PrivateKey)
}

func (c Jwt) PublicPath() string {
	return filepath.Join(c.CertsDir, c.PublicKey)
}
