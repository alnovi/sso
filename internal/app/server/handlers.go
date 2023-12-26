package server

import (
	"github.com/alnovi/sso/internal/transport/http/handler"
	"github.com/alnovi/sso/internal/transport/http/handler/api"
	"github.com/alnovi/sso/internal/transport/http/handler/web"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type ApiHandlers struct {
	auth    *api.Auth
	profile *api.Profile
}

type WebHandlers struct {
	auth    *web.Auth
	profile *web.Profile
	token   *web.Token
	home    *web.Home
}

type Handlers struct {
	web *WebHandlers
	api *ApiHandlers
	err *handler.Error
	doc echo.HandlerFunc
}

func newHandlers(app *App, uc *UseCases) (*Handlers, error) {
	apiHandlers := &ApiHandlers{
		auth:    api.NewAuth(uc.auth, uc.client),
		profile: api.NewProfile(uc.user),
	}

	webHandlers := &WebHandlers{
		auth:    web.NewAuth(uc.auth, uc.client),
		home:    web.NewHome(),
		profile: web.NewProfile(app.clients.profile, uc.token),
		token:   web.NewToken(uc.client, uc.token),
	}

	return &Handlers{
		web: webHandlers,
		api: apiHandlers,
		err: handler.NewError(),
		doc: echoSwagger.WrapHandler,
	}, nil
}
