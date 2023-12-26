package api

import (
	"net/http"

	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Profile struct {
	user usecase.User
}

func NewProfile(user usecase.User) *Profile {
	return &Profile{user: user}
}

func (h *Profile) UserInfo(c echo.Context) error {
	userId := c.Get(middleware.KeyUserId).(string)

	user, err := h.user.UserInfo(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	//TODO: информация о аватарке
	//TODO: информация о токенах
	//TODO: информация о доступных приложениях

	return c.JSON(http.StatusOK, response.User{
		UID:   user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
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
