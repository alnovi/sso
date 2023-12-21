package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) Profile(c echo.Context) error {
	return c.Render(http.StatusOK, "profile.html", nil)
}

func (h *ProfileHandler) UserInfo(c echo.Context) error {
	return nil
}

func (h *ProfileHandler) ChangeInfo(c echo.Context) error {
	return nil
}

func (h *ProfileHandler) ChangePassword(c echo.Context) error {
	return nil
}

func (h *ProfileHandler) Logout(c echo.Context) error {
	return nil
}
