package controller

import "github.com/labstack/echo/v4"

type ErrorController struct{}

func NewErrorHandler() *ErrorController {
	return &ErrorController{}
}

func (h *ErrorController) Handle(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
}
