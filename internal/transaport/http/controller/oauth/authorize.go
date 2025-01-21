package oauth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transaport/http/controller"
	"github.com/alnovi/sso/internal/transaport/http/request"
	"github.com/alnovi/sso/internal/transaport/http/response"
	"github.com/alnovi/sso/pkg/utils"
	"github.com/alnovi/sso/pkg/validator"
)

type AuthorizeController struct {
	*controller.BaseController
	oauth  *oauth.OAuth
	cookie *cookie.Cookie
}

func NewAuthorizeController(oauth *oauth.OAuth, cookie *cookie.Cookie) *AuthorizeController {
	return &AuthorizeController{oauth: oauth, cookie: cookie}
}

// Client       godoc
// @Id          CheckClient
// @Summary     Информация о клиенте
// @Description Проверка параметров клиента и информация о нем
// @Tags        OAuth
// @Accept      json
// @Produce     json
// @Param       response_type query string true "Response type" example(code)
// @Param       client_id     query string true "Client ID"     example(app_id)
// @Param       redirect_uri  query string true "Redirect URI"
// @Success 200
// @Router      /oauth/v1/client [get]
func (c *AuthorizeController) Client(e echo.Context) error {
	if e.QueryParam("response_type") != oauth.ResponseTypeCode {
		return echo.NewHTTPError(http.StatusBadRequest, "response type invalid")
	}

	client, err := c.oauth.Client(e.Request().Context(), e.QueryParam("client_id"), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
	}

	_, err = c.oauth.RedirectURL(client, e.QueryParam("redirect_uri"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "redirect url invalid").SetInternal(err)
	}

	return e.JSON(http.StatusOK, response.Client{
		Id:   client.Id,
		Name: client.Name,
	})
}

// Authorize    godoc
// @Id          Authorize
// @Summary     Аутентификация пользователя
// @Description Аутентификация пользователя
// @Tags        OAuth
// @Accept      json
// @Produce     json
// @Param       response_type query string            true  "Response type" example(code)
// @Param       client_id     query string            true  "Client ID"     example(app_id)
// @Param       redirect_uri  query string            true  "Redirect URI"
// @Param       state         query string            false "State"
// @Param       request       body  request.Authorize true "Логин и пароль пользователя"
// @Success 200 {object}            response.URL      "Ссылка для перехода"
// @Success 302
// @Router      /oauth/v1/authorize [post]
func (c *AuthorizeController) Authorize(e echo.Context) error {
	req := new(request.Authorize)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	if e.QueryParam("response_type") != oauth.ResponseTypeCode {
		return echo.NewHTTPError(http.StatusBadRequest, "response type invalid")
	}

	inp := oauth.InputAuthByCode{
		ClientId:    e.QueryParam("client_id"),
		RedirectUri: e.QueryParam("redirect_uri"),
		State:       e.QueryParam("state"),
		Login:       req.Login,
		Password:    req.Password,
		IP:          e.RealIP(),
		Agent:       e.Request().UserAgent(),
	}

	redirect, code, err := c.oauth.AuthorizeByCode(e.Request().Context(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrRedirectUriInvalid) {
			return echo.NewHTTPError(http.StatusBadRequest, "redirect url invalid").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrUserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "Логин не найден")
		}
		if errors.Is(err, oauth.ErrUserPasswordInvalid) {
			return validator.NewValidateErrorWithMessage("password", "Пароль не верный")
		}
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}

	if req.Remember {
		e.SetCookie(c.cookie.SessionId(*code.SessionId))
	}

	if utils.RequestIsAjax(e.Request()) {
		return e.JSON(http.StatusOK, response.URL{URL: redirect.String()})
	}

	return e.Redirect(http.StatusFound, redirect.String())
}

func (c *AuthorizeController) ApplyHTTP(g *echo.Group) error {
	g.GET("/client/", c.Client)
	g.POST("/authorize/", c.Authorize)
	return nil
}
