package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/jwt"
)

func Auth(jwt *jwt.JWT) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cls, err := jwt.ValidateAccessToken(c.Request().Header.Get(echo.HeaderAuthorization))
			if err != nil {
				return err
			}

			c.Set("user_id", cls.UserId())
			c.Set("client_id", cls.ClientId())
			c.Set("user_role", cls.UserRole())

			return next(c)
		}
	}
}
