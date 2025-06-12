package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/alnovi/gomon/utils"
	"github.com/alnovi/gomon/validator"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type AuthController struct {
	controller.BaseController
	oauth  *oauth.OAuth
	cookie *cookie.Cookie
}

func NewAuthController(oauth *oauth.OAuth, cookie *cookie.Cookie) *AuthController {
	return &AuthController{oauth: oauth, cookie: cookie}
}

func (c *AuthController) Form(e echo.Context) error {
	if session, err := e.Cookie(cookie.SessionId); err == nil {
		var redirectURI *url.URL

		inp := oauth.InputAuthorizeBySession{
			ClientId:     e.QueryParam("client_id"),
			ResponseType: e.QueryParam("response_type"),
			RedirectUri:  e.QueryParam("redirect_uri"),
			State:        e.QueryParam("state"),
			SessionId:    session.Value,
		}

		_, _, redirectURI, err = c.oauth.AuthorizeBySession(context.Background(), inp)
		if errors.Is(err, oauth.ErrSessionNotFound) {
			e.SetCookie(c.cookie.Remove(cookie.SessionId))
		}

		if redirectURI != nil {
			return e.Redirect(http.StatusFound, redirectURI.String())
		}
	}

	inp := oauth.InputAuthorizeParams{
		ClientId:     e.QueryParam("client_id"),
		ResponseType: e.QueryParam("response_type"),
		RedirectUri:  e.QueryParam("redirect_uri"),
	}

	client, err := c.oauth.AuthorizeCheckParams(context.Background(), inp)
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

func (c *AuthController) Authorize(e echo.Context) error {
	req := new(request.Authorize)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := oauth.InputAuthorizeByCode{
		ClientId:     e.QueryParam("client_id"),
		ResponseType: e.QueryParam("response_type"),
		RedirectUri:  utils.NormalizeURL(e.QueryParam("redirect_uri")),
		State:        e.QueryParam("state"),
		Login:        req.Login,
		Password:     req.Password,
		UserIP:       e.RealIP(),
		UserAgent:    e.Request().UserAgent(),
	}

	_, token, redirectURI, err := c.oauth.AuthorizeByCode(context.Background(), inp)
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
		if errors.Is(err, oauth.ErrUserNotFound) {
			return validator.NewValidateErrorWithMessage("login", "пользователь не найден")
		}
		if errors.Is(err, oauth.ErrInvalidUserPassword) {
			return validator.NewValidateErrorWithMessage("password", "пароль не верный")
		}
		return err
	}

	e.SetCookie(c.cookie.SessionId(*token.SessionId, req.Remember))

	if utils.RequestIsAjax(e.Request()) {
		return e.JSON(http.StatusOK, response.URL{URL: redirectURI.String()})
	}

	return e.Redirect(http.StatusFound, redirectURI.String())
}

func (c *AuthController) ApplyHTTP(g *echo.Group) {
	g.GET("/authorize/", c.Form)
	g.POST("/authorize/", c.Authorize)
}
