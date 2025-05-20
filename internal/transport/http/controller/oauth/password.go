package oauth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
)

type PasswordController struct {
	controller.BaseController
	oauth *oauth.OAuth
}

func NewPasswordController(oauth *oauth.OAuth) *PasswordController {
	return &PasswordController{oauth: oauth}
}

func (c *PasswordController) ForgotPassword(e echo.Context) error {
	req := new(request.ForgotPassword)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := oauth.InputForgotPassword{
		ClientId:    e.QueryParam("client_id"),
		RedirectUri: e.QueryParam("redirect_uri"),
		Query:       e.Request().URL.Query().Encode(),
		Login:       req.Login,
		IP:          e.RealIP(),
		Agent:       e.Request().UserAgent(),
	}

	if err := c.oauth.ForgotPassword(e.Request().Context(), inp); err != nil {
		if errors.Is(err, oauth.ErrUserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "пользователь не найден")
		}
		return err
	}

	return e.JSON(http.StatusOK, response.Message{
		Message: "Ссылка для смены пароля отправлена на электронную почту",
	})
}

func (c *PasswordController) ResetPassword(e echo.Context) error {
	req := new(request.ResetPassword)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := oauth.InputResetPassword{
		Hash:     e.QueryParam("hash"),
		Password: req.Password,
	}

	redirect, err := c.oauth.ResetPassword(e.Request().Context(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "token not found")
		}
		return err
	}

	if utils.RequestIsAjax(e.Request()) {
		return e.JSON(http.StatusOK, response.URL{URL: redirect.String()})
	}

	return e.Redirect(http.StatusFound, redirect.String())
}

func (c *PasswordController) ApplyHTTP(g *echo.Group) {
	g.POST("/forgot-password/", c.ForgotPassword)
	g.POST("/reset-password/", c.ResetPassword)
}
