package server

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/pkg/template"
	"github.com/alnovi/sso/internal/pkg/translate"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func NewHttpServer(app *App, m *Middlewares, h *Handlers) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.NewValidator()
	e.Renderer = template.NewHtmlRenderer(app.cfg.Path.Html)
	e.HTTPErrorHandler = httpErrorHandler

	e.File("/favicon.ico", fmt.Sprintf("%s/favicon.png", app.cfg.Path.Web))
	e.Static("/static/*", app.cfg.Path.Web)

	e.GET("/doc/*", h.doc)
	e.GET("/doc", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/doc/index.html")
	})

	e.GET("/oauth/signin", h.auth.AuthForm)
	e.POST("/oauth/signin", h.auth.SignIn)
	e.POST("/oauth/token", h.token.GenerateToken)

	return e, nil
}

func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		_ = c.JSON(echoHttpError.Code, response.Error{
			Message: echoHttpError.Message.(string),
		})
		return
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		_ = c.JSON(http.StatusUnprocessableEntity, response.ErrorValidate{
			Message:  translate.HttpStatusTextRU(http.StatusUnprocessableEntity),
			Validate: validateError.Fields,
		})
		return
	}

	if errors.Is(err, exception.AccessDenied) {
		_ = c.JSON(http.StatusForbidden, response.Error{
			Message: translate.HttpStatusTextRU(http.StatusForbidden),
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, response.Error{
		Message: translate.HttpStatusTextRU(http.StatusInternalServerError),
	})
}
