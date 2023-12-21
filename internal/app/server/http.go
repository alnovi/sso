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
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func NewHttpServer(app *App, m *Middlewares, h *Handlers) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.NewValidator()
	e.Renderer = template.NewHtmlRenderer(app.cfg.Path.Html)
	e.HTTPErrorHandler = httpErrorHandler

	e.Any("/", h.home.GoToAuth)

	e.GET("/oauth/signin", h.auth.Auth)
	e.POST("/oauth/signin", h.auth.SignIn)
	e.POST("/oauth/token", h.token.GenerateToken)

	e.File("/favicon.ico", fmt.Sprintf("%s/favicon.png", app.cfg.Path.Store))
	e.Static("/assets/*", app.cfg.Path.Assets)
	e.Static("/store/*", app.cfg.Path.Store)

	profile := e.Group("", m.profile)
	profile.GET("/profile", h.profile.Profile)

	e.GET("/doc/*", h.doc)

	return e, nil
}

func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	data := response.Error{
		Code:    http.StatusInternalServerError,
		Message: translate.HttpStatusTextRU(http.StatusInternalServerError),
	}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		data.Code = echoHttpError.Code
		data.Message = echoHttpError.Message.(string)

		if data.Message == http.StatusText(data.Code) {
			data.Message = translate.HttpStatusTextRU(data.Code)
		}
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		data.Code = http.StatusUnprocessableEntity
		data.Message = translate.HttpStatusTextRU(http.StatusUnprocessableEntity)
		data.Validate = validateError.Fields
	}

	if errors.Is(err, exception.ClientAccessDenied) {
		data.Code = http.StatusForbidden
		data.Message = translate.HttpStatusTextRU(http.StatusForbidden)
	}

	if errors.Is(err, exception.NotAuthorization) {
		_ = c.Redirect(http.StatusFound, "/oauth/signin")
		return
	}

	if utils.RequestIsJson(c.Request()) {
		_ = c.JSON(data.Code, data)
	} else {
		_ = c.Render(data.Code, "error.html", data)
	}
}
