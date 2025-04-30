package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Option func(e any)

func WithAccessExpiresAt(val time.Time) Option {
	return func(e any) {
		claims, ok := e.(*AccessClaims)
		if !ok {
			return
		}

		if val.IsZero() {
			return
		}

		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(val)
	}
}
