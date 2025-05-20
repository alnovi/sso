package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type TokenController struct {
	controller.BaseController
	oauth *oauth.OAuth
}

func NewTokenController(oauth *oauth.OAuth) *TokenController {
	return &TokenController{oauth: oauth}
}

func (c *TokenController) Token(e echo.Context) error {
	switch e.QueryParam("grant_type") {
	case oauth.GrantTypeAuthorizationCode:
		return c.tokenByCode(e)
	case oauth.GrantTypeRefreshToken:
		return c.tokenByRefresh(e)
	}
	return echo.NewHTTPError(http.StatusBadRequest, "grant_type is unsupported")
}

func (c *TokenController) tokenByCode(e echo.Context) error {
	inp := oauth.InputTokenByCode{
		ClientId:     e.QueryParam("client_id"),
		ClientSecret: e.QueryParam("client_secret"),
		Code:         e.QueryParam("code"),
	}

	access, refresh, err := c.oauth.TokenByCode(context.Background(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "token not found").SetInternal(err)
		}
		return err
	}

	return e.JSON(http.StatusOK, response.AccessToken{
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration,
	})
}

func (c *TokenController) tokenByRefresh(e echo.Context) error {
	inp := oauth.InputTokenByRefresh{
		ClientId:     e.QueryParam("client_id"),
		ClientSecret: e.QueryParam("client_secret"),
		Refresh:      e.QueryParam("refresh_token"),
	}

	access, refresh, err := c.oauth.TokenByRefresh(context.Background(), inp)
	if err != nil {
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "token not found").SetInternal(err)
		}
		return err
	}

	return e.JSON(http.StatusOK, response.AccessToken{
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration,
	})
}

func (c *TokenController) ApplyHTTP(g *echo.Group) {
	g.POST("/token/", c.Token)
}
