package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Session string `json:"session"`
	Client  string `json:"client"`
	User    string `json:"user"`
	Role    string `json:"role"`
}

func (c AccessClaims) SessionId() string {
	return c.Session
}

func (c AccessClaims) ClientId() string {
	return c.Client
}

func (c AccessClaims) UserId() string {
	return c.User
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
