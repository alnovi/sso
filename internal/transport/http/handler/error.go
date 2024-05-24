package handler

import (
	"errors"
	"net/http"

	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/helper"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
	BaseHandler
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Handle(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	data := response.Error{Code: http.StatusInternalServerError}

	var echoHttpError *echo.HTTPError
	if errors.As(err, &echoHttpError) {
		data.Code = echoHttpError.Code
		if http.StatusText(data.Code) != echoHttpError.Message.(string) {
			data.Message = echoHttpError.Message.(string)
		}
	}

	var validateError *validator.ValidateError
	if errors.As(err, &validateError) {
		data.Code = http.StatusUnprocessableEntity
		data.Message = h.StatusText(http.StatusUnprocessableEntity)
		data.Validate = validateError.Fields
	}

	if data.Message == "" {
		data.Message = h.StatusText(data.Code)
	}

	if helper.RequestIsJson(c.Request()) {
		_ = c.JSON(data.Code, data)
	} else {
		_ = c.Render(data.Code, "error.html", data)
	}
}
