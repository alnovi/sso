package certs

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type JWK struct {
	Kty    string   `json:"kty"`
	Alg    string   `json:"alg"`
	Use    string   `json:"use,omitempty"`
	KeyOps []string `json:"key_ops,omitempty"`
	N      string   `json:"n"`
	E      string   `json:"e"`
}

func NewJwk(key *rsa.PublicKey) *JWK {
	return &JWK{
		Kty:    "RSA",
		Alg:    "RS256",
		Use:    "sig",
		KeyOps: []string{"verify"},
		N:      base64.URLEncoding.EncodeToString(key.N.Bytes()),
		E:      base64.URLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}
