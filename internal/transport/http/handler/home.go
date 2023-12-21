package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) GoToAuth(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/oauth/signin")
}
