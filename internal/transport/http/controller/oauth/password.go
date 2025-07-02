package oauth

import (
	"errors"
	"net/http"

	"github.com/alnovi/gomon/utils"
	"github.com/alnovi/gomon/validator"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type PasswordController struct {
	controller.BaseController
	oauth *oauth.OAuth
}

func NewPasswordController(oauth *oauth.OAuth) *PasswordController {
	return &PasswordController{oauth: oauth}
}

// FormForgot   godoc
// @Id          OAuthFormForgot
// @Summary     Форма "Забыли пароль"
// @Description Форма получения токена для восстановления пароля
// @Tags        OAuth
// @Accept      html
// @Produce     html
// @Param       client_id query string true "идентификатор клиента"
// @Param       response_type query string true "тип запроса"
// @Param       redirect_uri query string true "адрес клиента (callback)"
// @Param       state query string false "состояние"
// @Success 200
// @Failure 400
// @Router      /oauth/forgot-password [get]
func (c *PasswordController) FormForgot(e echo.Context) error {
	inp := oauth.InputAuthorizeParams{
		ClientId:     e.QueryParam("client_id"),
		ResponseType: e.QueryParam("response_type"),
		RedirectUri:  e.QueryParam("redirect_uri"),
	}

	client, err := c.oauth.AuthorizeCheckParams(e.Request().Context(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrInvalidResponseType) {
			return echo.NewHTTPError(http.StatusBadRequest, "Не валидный response-type").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Клиент не найден").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrInvalidRedirectUri) {
			return echo.NewHTTPError(http.StatusBadRequest, "Не валидный redirect-uri").SetInternal(err)
		}
		return err
	}

	resp := echo.Map{
		"Version": config.Version,
		"Query":   e.Request().URL.RawQuery,
		"Name":    client.Name,
		"Icon":    client.Icon,
	}

	return e.Render(http.StatusOK, "auth.html", resp)
}

// FormReset    godoc
// @Id          OAuthFormReset
// @Summary     Форма восстановления пароля
// @Description Форма для восстановления доступа
// @Tags        OAuth
// @Accept      html
// @Produce     html
// @Param       hash query string true "токен восстановления доступа"
// @Success 200
// @Failure 400
// @Router      /oauth/reset-password [get]
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

// ForgotPassword godoc
// @Id            OAuthForgotPassword
// @Summary       Отправка токена сброса пароля
// @Description   Отправка токена сброса пароля на почту пользователя
// @Tags          OAuth
// @Accept        json
// @Produce       json
// @Param         request body request.ForgotPassword true "Логин пользователя"
// @Success 200
// @Failure 400
// @Failure 422
// @Router        /oauth/forgot-password [post]
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

// ResetPassword godoc
// @Id           OAuthResetPassword
// @Summary      Сброс пароля
// @Description  Задать новый пароль пользователя
// @Tags         OAuth
// @Accept       json
// @Produce      json
// @Param        request body request.ResetPassword true "Новый пароль пользователя"
// @Success 200
// @Success 302
// @Failure 400
// @Failure 422
// @Router       /oauth/reset-password [post]
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
	g.GET("/forgot-password/", c.FormForgot)
	g.POST("/forgot-password/", c.ForgotPassword)
	g.GET("/reset-password/", c.FormReset)
	g.POST("/reset-password/", c.ResetPassword)
}
