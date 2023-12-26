package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/pkg/cookies"
	"github.com/alnovi/sso/internal/service"
	"github.com/labstack/echo/v4"
)

const (
	KeyToken    = "token"
	KeyUserId   = "user_id"
	KeyClientId = "client_id"
)

func AuthProfile(profile *entity.Client, ts service.Token) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			hash, err := tokenHash(c, cookies.NameProfileAccess)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
			}

			token, err := ts.Validate(ctx, profile.Id, profile.Secret, entity.TokenClassAccess, hash)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
			}

			c.Set(KeyToken, *token)
			c.Set(KeyUserId, *token.UserId)
			c.Set(KeyClientId, *token.ClientId)

			return next(c)
		}
	}
}

func tokenHash(c echo.Context, cookieName string) (string, error) {
	if tokenHeader := c.Request().Header.Get("Authorization"); len(tokenHeader) > 0 {
		return strings.TrimPrefix(tokenHeader, "Bearer "), nil
	}

	if tokenQuery := c.QueryParam("token"); len(tokenQuery) > 0 {
		return tokenQuery, nil
	}

	if tokenCookie, err := c.Cookie(cookieName); err == nil {
		return tokenCookie.Value, nil
	}

	return "", errors.New("token not found")
}
