package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func (c AccessClaims) ClientId() string {
	return c.RegisteredClaims.Issuer
}

func (c AccessClaims) UserId() string {
	return c.RegisteredClaims.Subject
}

func (c AccessClaims) UserRole() string {
	return c.Role
}

func (c AccessClaims) NotBefore() time.Time {
	return c.RegisteredClaims.NotBefore.Time
}

func (c AccessClaims) ExpiresAt() time.Time {
	return c.RegisteredClaims.ExpiresAt.Time
}
