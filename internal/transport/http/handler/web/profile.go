package web

import (
	"net/http"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/pkg/cookies"
	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Profile struct {
	profile *entity.Client
	token   usecase.Token
}

func NewProfile(profile *entity.Client, token usecase.Token) *Profile {
	return &Profile{profile: profile, token: token}
}

func (h *Profile) Profile(c echo.Context) error {
	token, ok := c.Get(middleware.KeyToken).(entity.Token)
	if !ok || *token.ClientId != h.profile.Id {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	return c.Render(http.StatusOK, "profile.html", echo.Map{
		"AppName":     h.profile.Name,
		"AccessToken": token.Hash,
	})
}

func (h *Profile) ProfileCallback(c echo.Context) error {
	code := c.QueryParam("code")

	dtoAccess := dto.AccessToken{
		Client:    *h.profile,
		Code:      code,
		GrantType: dto.GrantTypeAuthorizationCode,
	}

	access, refresh, err := h.token.AccessAndRefreshToken(c.Request().Context(), dtoAccess)
	if exception.Is(err) {
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	c.SetCookie(cookies.ProfileAccess(access.Hash))
	c.SetCookie(cookies.ProfileRefresh(refresh.Hash))

	return c.Redirect(http.StatusFound, "/profile")
}
