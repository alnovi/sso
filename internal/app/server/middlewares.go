package server

import (
	"context"

	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/labstack/echo/v4"
)

type Middlewares struct {
	profile echo.MiddlewareFunc
}

func newMiddlewares(_ *App, s *Services) (*Middlewares, error) {
	var err error

	clientProfile, err := s.client.GetProfileClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Middlewares{
		profile: middleware.AuthProfile(clientProfile.Id, clientProfile.Secret, s.token),
	}, nil
}
