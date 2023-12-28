package api

import (
	"net/http"

	"github.com/alnovi/sso/internal/transport/http/middleware"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Profile struct {
	profile usecase.Profile
}

func NewProfile(profile usecase.Profile) *Profile {
	return &Profile{profile: profile}
}

func (h *Profile) UserInfo(c echo.Context) error {
	userId := c.Get(middleware.KeyUserId).(string)

	user, tokens, clients, err := h.profile.Profile(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	res := response.Profile{
		Id:      user.Id,
		Name:    user.Name,
		Image:   user.Image,
		Email:   user.Email,
		Tokens:  make([]response.ProfileToken, 0),
		Clients: make([]response.ProfileClient, 0),
	}

	for _, token := range tokens {
		meta := new(response.ProfileTokenMeta)

		if token.Meta != nil {
			meta.IP = token.Meta.IP
			meta.Agent = token.Meta.Agent
		}

		res.Tokens = append(res.Tokens, response.ProfileToken{
			Id:        token.Id,
			Class:     token.Class,
			Meta:      meta,
			CreatedAt: token.CreatedAt,
			UpdatedAt: token.UpdatedAt,
		})
	}

	for _, client := range clients {
		res.Clients = append(res.Clients, response.ProfileClient{
			Id:          client.Id,
			Name:        client.Name,
			Description: client.Description,
			Logo:        client.Logo,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Profile) ChangeInfo(c echo.Context) error {
	return nil
}

func (h *Profile) ChangePassword(c echo.Context) error {
	return nil
}

func (h *Profile) Logout(c echo.Context) error {
	return nil
}
