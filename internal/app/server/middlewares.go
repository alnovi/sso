package server

import (
	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
)

type Middlewares struct {
	logger  echo.MiddlewareFunc
	profile echo.MiddlewareFunc
}

func newMiddlewares(app *App, s *Services) (*Middlewares, error) {
	return &Middlewares{
		logger:  middleware.RequestLogger(app.log),
		profile: middleware.AuthProfile(app.clients.profile, s.token),
	}, nil
}
