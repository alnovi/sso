package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/jwt"
	"github.com/alnovi/sso/internal/transaport/http/response"
	"github.com/alnovi/sso/pkg/server"
	"github.com/alnovi/sso/pkg/validator"
)

type ErrorController struct{}

func NewErrorHandler() *ErrorController {
	return &ErrorController{}
}

func (h *ErrorController) Handle(err error, c echo.Context) {
	if c.Response().Committed {
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

	if errors.Is(err, jwt.ErrUnauthenticated) {
		data.Code = http.StatusUnauthorized
		data.Error = server.StatusText(http.StatusUnauthorized)
	}

	if data.Error == "" {
		data.Error = server.StatusText(data.Code)
	}

	_ = c.JSON(data.Code, data)
}
