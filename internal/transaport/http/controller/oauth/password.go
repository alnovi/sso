package oauth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transaport/http/controller"
	"github.com/alnovi/sso/internal/transaport/http/request"
	"github.com/alnovi/sso/internal/transaport/http/response"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
)

type PasswordController struct {
	*controller.BaseController
	oauth *oauth.OAuth
}

func NewPasswordController(oauth *oauth.OAuth) *PasswordController {
	return &PasswordController{oauth: oauth}
}

// ForgotPassword godoc
// @Id            ForgotPassword
// @Summary       Забыли пароль
// @Description   Отправка ссылки для смены пароля
// @Tags          OAuth-password
// @Accept        json
// @Produce       json
// @Param         response_type query string                 true  "Response type" example(code)
// @Param         client_id     query string                 true  "Client ID"     example(app_id)
// @Param         redirect_uri  query string                 true  "Redirect URI"
// @Param         state         query string                 false "State"
// @Param         request       body  request.ForgotPassword true  "Логин пользователя"
// @Success 200   {object}            response.Message             "Сообщение для пользователя"
// @Router        /v1/oauth/forgot-password [post]
func (c *PasswordController) ForgotPassword(e echo.Context) error {
	req := new(request.ForgotPassword)

	if _, err := c.oauth.ResponseType(e.QueryParam("response_type")); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := oauth.InputForgotPassword{
		ClientId:    e.QueryParam("client_id"),
		RedirectUri: e.QueryParam("redirect_uri"),
		Query:       e.Request().URL.Query().Encode(),
		Login:       req.Login,
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

// ResetPassword godoc
// @Id           ResetPassword
// @Summary      Смена пароля
// @Description  Изменение пароля пользователя
// @Tags         OAuth-password
// @Accept       json
// @Produce      json
// @Param        hash     query string                true "Разовый токен" example(secret)
// @Param        request  body  request.ResetPassword true "Новый пароль пользователя"
// @Success 200  {object}       response.URL               "Ссылка для перехода"
// @Router       /v1/oauth/reset-password [post]
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

func (c *PasswordController) ApplyHTTP(g *echo.Group) error {
	g.POST("/forgot-password/", c.ForgotPassword)
	g.POST("/reset-password/", c.ResetPassword)
	return nil
}
