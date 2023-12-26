package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Home struct {
}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) GoToAuth(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/oauth/signin")
}
