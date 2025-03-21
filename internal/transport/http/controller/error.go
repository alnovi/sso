package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/server"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
)

type ErrorController struct{}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (c *ErrorController) Handle(err error, e echo.Context) {
	if e.Response().Committed {
		return
	}

	data := response.Error{Code: http.StatusInternalServerError}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		data.Code = echoHttpError.Code
		data.Error = echoHttpError.Message.(string)

		if data.Error == http.StatusText(data.Code) {
			data.Error = server.StatusText(data.Code)
		}

		if data.Code == http.StatusTooManyRequests {
			data.Error = server.StatusText(data.Code)
		}
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		data.Code = http.StatusUnprocessableEntity
		data.Error = server.StatusText(http.StatusUnprocessableEntity)
		data.Validate = validateError.Fields
	}

	if errors.Is(err, repository.ErrNoResult) {
		data.Code = http.StatusNotFound
	}

	if errors.Is(err, oauth.ErrUnauthorized) || errors.Is(err, echo.ErrUnauthorized) {
		data.Code = http.StatusUnauthorized
	}

	if errors.Is(err, oauth.ErrForbidden) || errors.Is(err, echo.ErrForbidden) {
		data.Code = http.StatusForbidden
	}

	if data.Error == "" {
		data.Error = server.StatusText(data.Code)
	}

	_ = c.Render(e, data)
}

func (c *ErrorController) Render(e echo.Context, data response.Error) error {
	if utils.RequestIsJson(e.Request()) {
		return e.JSON(data.Code, data)
	}
	return e.Render(data.Code, "error.html", data)
}
