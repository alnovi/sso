package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/service/oauth"
)

type ProfileController struct {
	oauth  *oauth.OAuth
	cookie *cookie.Cookie
}

func NewProfileController(oauth *oauth.OAuth, cookie *cookie.Cookie) *ProfileController {
	return &ProfileController{oauth: oauth, cookie: cookie}
}

// Logout       godoc
// @Id          Logout
// @Summary     Выход
// @Description Удаление сессии пользователя
// @Tags        OAuth
// @Accept      json
// @Produce     json
// @Success 200
// @Router      /oauth/v1/logout [post]
func (c *ProfileController) Logout(e echo.Context) error {
	if session, _ := e.Cookie(cookie.SessionId); session != nil {
		_ = c.oauth.RemoveSession(e.Request().Context(), session.Value)
	}

	e.SetCookie(c.cookie.Remove(cookie.SessionId))

	return e.NoContent(http.StatusOK)
}

func (c *ProfileController) ApplyHTTP(g *echo.Group) error {
	g.POST("/logout/", c.Logout)
	return nil
}
