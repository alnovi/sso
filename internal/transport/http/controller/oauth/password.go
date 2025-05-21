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

func (c *PasswordController) FormReset(e echo.Context) error {
	token, client, err := c.oauth.ValidateForgotToken(e.Request().Context(), e.QueryParam("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Токен не найден").SetInternal(err)
	}

	resp := echo.Map{
		"Query": token.Payload.Query(),
		"Name":  client.Name,
		"Icon":  client.Icon,
	}

	return e.Render(http.StatusOK, "auth.html", resp)
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
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Клиент не найден").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrInvalidRedirectUri) {
			return echo.NewHTTPError(http.StatusBadRequest, "Не валидный redirect-uri").SetInternal(err)
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
		Hash:     req.Token,
		Password: req.Password,
	}

	redirect, err := c.oauth.ResetPassword(e.Request().Context(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Токен не найден").SetInternal(err)
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
	g.GET("/reset-password/", c.FormReset)
	g.POST("/reset-password/", c.ResetPassword)
}
