package server

import (
	"github.com/alnovi/sso/internal/transport/http/handler"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	doc   echo.HandlerFunc
	auth  *handler.AuthHandler
	token *handler.TokenHandler
}

func newHandlers(_ *App, uc *UseCases) (*Handlers, error) {
	return &Handlers{
		doc:   echoSwagger.WrapHandler,
		auth:  handler.NewAuthHandler(uc.auth, uc.client),
		token: handler.NewTokenHandler(),
	}, nil
}
