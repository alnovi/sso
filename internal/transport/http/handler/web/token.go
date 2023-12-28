package web

import (
	"errors"
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Token struct {
	client usecase.Client
	token  usecase.Token
}

func NewToken(client usecase.Client, token usecase.Token) *Token {
	return &Token{client: client, token: token}
}

func (h *Token) GenerateToken(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	clientId := c.QueryParam("client_id")
	clientSecret := c.QueryParam("client_secret")
	codeToken := c.QueryParam("code")
	refreshToken := c.QueryParam("refresh_token")
	grantType := c.QueryParam("grant_type")

	dtoClient := dto.ClientForToken{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	client, err := h.client.ClientForToken(ctx, dtoClient)
	if err != nil {
		if exception.Is(err) {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	dtoToken := dto.AccessToken{
		Client:    *client,
		Code:      codeToken,
		Refresh:   refreshToken,
		GrantType: grantType,
	}

	access, refresh, err := h.token.AccessAndRefreshToken(ctx, dtoToken)
	if err != nil {
		if errors.Is(err, exception.ErrClientAccessDenied) {
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}
		if exception.Is(err) {
			return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	user := response.User{
		UID:   access.User.Id,
		Name:  access.User.Name,
		Image: access.User.Image,
		Email: access.User.Email,
	}

	return c.JSON(http.StatusOK, response.AccessToken{
		TokenType:    "bearer",
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration.Unix(),
		Info:         user,
	})
}
