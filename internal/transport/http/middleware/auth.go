package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/controller"
)

func Auth(auth *oauth.OAuth, cook *cookie.Cookie, clientId, clientSecret string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			ctx := e.Request().Context()

			accessToken := e.Request().Header.Get("Authorization")
			refreshToken := e.Request().Header.Get("Refresh-Token")

			if sessionToken, err := e.Cookie(cookie.NameAccessToken(clientId)); err == nil {
				accessToken = sessionToken.Value
			}

			if sessionToken, err := e.Cookie(cookie.NameRefreshToken(clientId)); err == nil {
				refreshToken = sessionToken.Value
			}

			if accessToken == "" {
				return echo.ErrUnauthorized
			}

			claims, err := auth.ValidateAccessToken(ctx, accessToken)
			if err != nil {
				var access *entity.Token
				var refresh *entity.Token

				if refreshToken == "" {
					return echo.ErrUnauthorized
				}

				inp := oauth.InputTokenByRefresh{
					ClientId:     clientId,
					ClientSecret: clientSecret,
					Refresh:      refreshToken,
				}

				refresh, err = auth.ValidateRefreshToken(ctx, refreshToken)
				if err != nil {
					return fmt.Errorf("%w: %s", echo.ErrUnauthorized, err)
				}

				if *refresh.ClientId != clientId {
					return fmt.Errorf("%w: client not attempted", echo.ErrUnauthorized)
				}

				access, refresh, err = auth.TokenByRefresh(ctx, inp)
				if err != nil {
					return fmt.Errorf("%w: %s", echo.ErrUnauthorized, err)
				}

				claims, err = auth.ValidateAccessToken(ctx, access.Hash)
				if err != nil {
					return fmt.Errorf("%w: %s", echo.ErrUnauthorized, err)
				}

				e.SetCookie(cook.AccessToken(access))
				e.SetCookie(cook.RefreshToken(refresh))
			}

			if claims.ClientId() != clientId {
				return fmt.Errorf("%w: client not attempted", echo.ErrUnauthorized)
			}

			e.Set(controller.CtxSessionId, claims.SessionId())
			e.Set(controller.CtxClientId, claims.ClientId())
			e.Set(controller.CtxUserId, claims.UserId())
			e.Set(controller.CtxUserRole, claims.UserRole())

			return next(e)
		}
	}
}
