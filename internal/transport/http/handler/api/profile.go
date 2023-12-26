package api

import (
	"github.com/labstack/echo/v4"
)

type Profile struct {
}

func NewProfile() *Profile {
	return &Profile{}
}

func (h *Profile) UserInfo(c echo.Context) error {
	return nil
}

func (h *Profile) ChangeInfo(c echo.Context) error {
	return nil
}

func (h *Profile) ChangePassword(c echo.Context) error {
	return nil
}

func (h *Profile) Logout(c echo.Context) error {
	return nil
}
