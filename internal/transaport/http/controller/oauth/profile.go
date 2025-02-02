package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/service/users"
	"github.com/alnovi/sso/internal/transaport/http/controller"
	"github.com/alnovi/sso/internal/transaport/http/response"
)

type ProfileController struct {
	*controller.BaseController
	oauth  *oauth.OAuth
	user   *users.User
	cookie *cookie.Cookie
}

func NewProfileController(oauth *oauth.OAuth, cookie *cookie.Cookie, user *users.User) *ProfileController {
	return &ProfileController{oauth: oauth, cookie: cookie, user: user}
}

// Profile      godoc
// @Id          Profile
// @Summary     Профиль пользователя
// @Description Профиль пользователя
// @Tags        OAuth profile
// @Accept      json
// @Produce     json
// @Success 200 {object} response.Profile "Профиль пользователя"
// @Router      /oauth/v1/profile [get]
func (c *ProfileController) Profile(e echo.Context) error {
	ctx := e.Request().Context()
	userId := c.MustUserId(e)

	user, err := c.user.UserById(ctx, userId)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.Profile{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// Logout       godoc
// @Id          Logout
// @Summary     Выход
// @Description Удаление сессии пользователя
// @Tags        OAuth profile
// @Accept      json
// @Produce     json
// @Success 200
// @Router      /oauth/v1/profile/logout [post]
func (c *ProfileController) Logout(e echo.Context) error {
	if session, _ := e.Cookie(cookie.SessionId); session != nil {
		_ = c.oauth.RemoveSession(e.Request().Context(), session.Value)
	}

	e.SetCookie(c.cookie.Remove(cookie.SessionId))

	return e.NoContent(http.StatusOK)
}

func (c *ProfileController) ApplyHTTP(g *echo.Group) error {
	g.POST("/profile/logout/", c.Logout)
	return nil
}
