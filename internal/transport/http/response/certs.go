package response

import "github.com/alnovi/sso/internal/service/certs"

type JWK struct {
	Kty    string   `json:"kty" example:"RSA"`
	Alg    string   `json:"alg" example:"RS256"`
	Use    string   `json:"use,omitempty" example:"sig"`
	KeyOps []string `json:"key_ops,omitempty" example:"verify"`
	N      string   `json:"n" example:"4TBai9lUV9qU0LqD37qpNTpJ1QuIn_2syDc9clGXKvf5lk6ESWaouvAT"`
	E      string   `json:"e" example:"AQAB"`
}

func NewJWK(jwk *certs.JWK) *JWK {
	return &JWK{
		Kty:    jwk.Kty,
		Alg:    jwk.Alg,
		Use:    jwk.Use,
		KeyOps: jwk.KeyOps,
		N:      jwk.N,
		E:      jwk.E,
	}
}
