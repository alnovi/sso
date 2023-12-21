package server

import (
	"github.com/alnovi/sso/internal/transport/http/handler"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	auth    *handler.AuthHandler
	home    *handler.HomeHandler
	profile *handler.ProfileHandler
	token   *handler.TokenHandler
	doc     echo.HandlerFunc
}

func newHandlers(_ *App, uc *UseCases) (*Handlers, error) {
	return &Handlers{
		auth:    handler.NewAuthHandler(uc.auth, uc.client),
		home:    handler.NewHomeHandler(),
		profile: handler.NewProfileHandler(),
		token:   handler.NewTokenHandler(uc.client, uc.token),
		doc:     echoSwagger.WrapHandler,
	}, nil
}
