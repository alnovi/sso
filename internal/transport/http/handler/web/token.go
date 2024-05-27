package web

import (
	"context"
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/labstack/echo/v4"
)

type tokenUseCase interface {
	AccessTokenByCode(ctx context.Context, inp dto.AccessTokenByCode) (*entity.Token, *entity.Token, error)
	AccessTokenByRefresh(ctx context.Context, inp dto.AccessTokenByRefresh) (*entity.Token, *entity.Token, error)
}

type TokenHandler struct {
	uc tokenUseCase
}

func NewTokenHandler(uc tokenUseCase) *TokenHandler {
	return &TokenHandler{uc: uc}
}

func (h *TokenHandler) Token(c echo.Context) error {
	var err error
	var access *entity.Token
	var refresh *entity.Token
	var ctx = c.Request().Context()

	switch c.QueryParam("grant_type") {
	case dto.GrantTypeCode:
		access, refresh, err = h.uc.AccessTokenByCode(ctx, dto.AccessTokenByCode{
			ClientID:     c.QueryParam("client_id"),
			ClientSecret: c.QueryParam("client_secret"),
			CodeHash:     c.QueryParam("code"),
		})
	case dto.GrantTypeRefresh:
		access, refresh, err = h.uc.AccessTokenByRefresh(ctx, dto.AccessTokenByRefresh{
			ClientID:     c.QueryParam("client_id"),
			ClientSecret: c.QueryParam("client_secret"),
			RefreshHash:  c.QueryParam("refresh_token"),
		})
	default:
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrUnsupportedGrantType)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	return c.JSON(http.StatusOK, response.AccessToken{
		TokenType:    "Bearer",
		AccessToken:  access.Hash,
		RefreshToken: refresh.Hash,
		ExpiresIn:    access.Expiration.Unix(),
	})
}

func (h *TokenHandler) Route(e *echo.Group) {
	e.POST("/oauth/token/", h.Token)
}
