package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/profile"
	"github.com/alnovi/sso/internal/transport/http/controller"
)

func AuthBySession(profile *profile.UserProfile) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			sessionId := e.QueryParam(cookie.SessionId)
			userAgent := e.Request().UserAgent()

			if cookieSession, err := e.Cookie(cookie.SessionId); err == nil {
				sessionId = cookieSession.Value
			}

			session, err := profile.SessionByIdAndAgent(e.Request().Context(), sessionId, userAgent)
			if err != nil {
				return fmt.Errorf("%w: %s", echo.ErrUnauthorized, err)
			}

			e.Set(controller.CtxSessionId, session.Id)
			e.Set(controller.CtxUserId, session.UserId)

			return next(e)
		}
	}
}
