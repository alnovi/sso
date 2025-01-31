package oauth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transaport/http/controller"
	"github.com/alnovi/sso/internal/transaport/http/response"
)

type TokenController struct {
	*controller.BaseController
	oauth *oauth.OAuth
}

func NewTokenController(oauth *oauth.OAuth) *TokenController {
	return &TokenController{oauth: oauth}
}

// Token  godoc
// @Id          Token
// @Summary     Токен доступа
// @Description Получение токена доступа по code или refresh токену
// @Tags        OAuth
// @Accept      json
// @Produce     json
// @Param       grant_type    query string true  "Grant type"    Enums(authorization_code, refresh_token)
// @Param       client_id     query string true  "Client ID"     example(app_id)
// @Param       client_secret query string true  "Client secret" example(secret)
// @Param       code          query string false "Code token"    example(secret)
// @Param       refresh_token query string false "Refresh token" example(secret)
// @Success 200 {object}      response.AccessToken "Токен доступа"
// @Router      /v1/oauth/token [post]
func (c *TokenController) Token(e echo.Context) error {
	switch e.QueryParam("grant_type") {
	case oauth.GrantTypeAuthorizationCode:
		return c.tokenByCode(e)
	case oauth.GrantTypeRefreshToken:
		return c.tokenByRefresh(e)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid grant_type")
	}
}

func (c *TokenController) tokenByCode(e echo.Context) error {
	ctx := e.Request().Context()

	inp := oauth.InputTokenByCode{
		ClientId:     e.QueryParam("client_id"),
		ClientSecret: e.QueryParam("client_secret"),
		Code:         e.QueryParam("code"),
	}

	access, refresh, err := c.oauth.TokenByCode(ctx, inp)
	if err != nil {
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "code incorrect").SetInternal(err)
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
	ctx := e.Request().Context()

	inp := oauth.InputTokenByRefresh{
		ClientId:     e.QueryParam("client_id"),
		ClientSecret: e.QueryParam("client_secret"),
		Refresh:      e.QueryParam("refresh_token"),
	}

	access, refresh, err := c.oauth.TokenByRefresh(ctx, inp)
	if err != nil {
		if errors.Is(err, oauth.ErrClientNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "client not found").SetInternal(err)
		}
		if errors.Is(err, oauth.ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "token incorrect").SetInternal(err)
		}
		return err
	}

	return e.JSON(http.StatusOK, response.AccessToken{
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration,
	})
}

func (c *TokenController) ApplyHTTP(g *echo.Group) error {
	g.POST("/token/", c.Token)
	return nil
}
