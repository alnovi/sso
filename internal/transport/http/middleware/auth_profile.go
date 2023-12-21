package middleware

import (
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/service"
	"github.com/labstack/echo/v4"
)

func AuthProfile(clientId, clientSecret string, ts service.Token) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			hash, err := getTokenHash(c)
			if err != nil {
				return err
			}

			token, err := ts.Validate(ctx, clientId, clientSecret, entity.TokenClassAccess, hash)
			if err != nil {
				return exception.NotAuthorization
			}

			c.Set("userId", token.UserId)
			c.Set("clientId", token.ClientId)

			return next(c)
		}
	}
}
