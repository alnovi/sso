package config

type Jwt struct {
	PrivateKey string `env:"PRIVATE_KEY"`
	PublicKey  string `env:"PUBLIC_KEY"`
}
