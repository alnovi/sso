package server

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/pkg/template"
	"github.com/alnovi/sso/internal/pkg/translate"
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

	code := http.StatusInternalServerError
	message := ""
	data := echo.Map{}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		code = echoHttpError.Code
		message = echoHttpError.Message.(string)
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		code = http.StatusUnprocessableEntity
		data = echo.Map{
			"validate": validateError.Fields,
		}
	}

	if errors.Is(err, exception.AccessDenied) {
		code = http.StatusForbidden
	}

	if message == "" || message == http.StatusText(code) {
		message = translate.HttpStatusTextRU(code)
	}

	data["code"] = code
	data["message"] = message

	_ = c.JSON(code, data)
}
