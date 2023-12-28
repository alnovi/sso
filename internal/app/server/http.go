package server

import (
	"fmt"

	"github.com/alnovi/sso/internal/transport/http/render"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func NewHttpServer(app *App, m *Middlewares, h *Handlers) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.NewValidator()
	e.Renderer = render.New(app.cfg.Path.Html)
	e.HTTPErrorHandler = h.err.ErrorHandle

	e.Use(m.logger)

	initWebRoutes(app, e, m, h)
	initApiRoutes(app, e, m, h)
	initOtherRoutes(app, e, m, h)

	return e, nil
}

func initWebRoutes(_ *App, e *echo.Echo, m *Middlewares, h *Handlers) {
	e.Any("/", h.web.home.GoToAuth)

	e.GET("/oauth/signin", h.web.auth.Auth)
	e.POST("/oauth/signin", h.web.auth.SignIn)
	e.POST("/oauth/token", h.web.token.GenerateToken)

	e.GET("/profile", h.web.profile.Profile, m.profile)
	e.GET("/profile/callback", h.web.profile.ProfileCallback)
}

func initApiRoutes(_ *App, e *echo.Echo, m *Middlewares, h *Handlers) {
	g := e.Group("/api")

	g.POST("/oauth/signin", h.api.auth.SignIn)

	g.GET("/profile", h.api.profile.UserInfo, m.profile)
	g.PUT("/profile", h.api.profile.ChangeInfo, m.profile)
	g.PUT("/profile/password", h.api.profile.ChangePassword, m.profile)
}

func initOtherRoutes(app *App, e *echo.Echo, m *Middlewares, h *Handlers) {
	e.File("/favicon.ico", fmt.Sprintf("%s/favicon.png", app.cfg.Path.Store))
	e.Static("/assets/*", app.cfg.Path.Assets)
	e.Static("/store/*", app.cfg.Path.Store)
}
