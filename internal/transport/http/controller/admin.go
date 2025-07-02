package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/service/admin"
	"github.com/alnovi/sso/internal/service/cookie"
)

type AdminController struct {
	BaseController
	admin  *admin.Admin
	cookie *cookie.Cookie
	token  echo.MiddlewareFunc
}

func NewAdminController(admin *admin.Admin, cookie *cookie.Cookie, token echo.MiddlewareFunc) *AdminController {
	return &AdminController{admin: admin, cookie: cookie, token: token}
}

// Home         godoc
// @Id          AdminHome
// @Summary     Раздел администратора
// @Description Раздел администратора. Управление клиентами, пользователями и т.д.
// @Tags        Admin
// @Accept      html
// @Produce     html
// @Security    JWT-Access
// @Success 200
// @Failure 302
// @Router      /admin [get]
func (c *AdminController) Home(e echo.Context) error {
	if _, ok := c.UserId(e); !ok {
		authorizeURL, err := c.admin.AuthorizeURI(context.Background())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
		}
		return e.Redirect(http.StatusFound, authorizeURL)
	}

	return e.Render(http.StatusOK, "admin.html", echo.Map{"Version": config.Version})
}

// Callback     godoc
// @Id          AdminCallback
// @Summary     Аутентификация пользователя
// @Description Аутентификация пользователя по коду. Вызывается автоматически после ввода логина и пароля.
// @Tags        Admin
// @Accept      html
// @Produce     html
// @Security    JWT-Access
// @Success     302
// @Header      302 {string}  Location  "/admin"
// @Router      /admin/callback [get]
func (c *AdminController) Callback(e echo.Context) error {
	access, refresh, err := c.admin.TokenByCode(context.Background(), e.QueryParam("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	e.SetCookie(c.cookie.AccessToken(access))
	e.SetCookie(c.cookie.RefreshToken(refresh))

	return e.Redirect(http.StatusFound, "/admin")
}

// Logout       godoc
// @Id          AdminLogout
// @Summary     Выход из системы
// @Description Выход из раздела администратора
// @Tags        Admin
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Success 200
// @Router      /admin/logout [post]
func (c *AdminController) Logout(e echo.Context) error {
	if sessionId, ok := c.SessionId(e); ok {
		_ = c.admin.Logout(context.Background(), sessionId)
	}
	e.SetCookie(c.cookie.Remove(cookie.SessionId))
	e.SetCookie(c.cookie.Remove(cookie.NameAccessToken(c.admin.ClientId())))
	e.SetCookie(c.cookie.Remove(cookie.NameRefreshToken(c.admin.ClientId())))
	return e.NoContent(http.StatusOK)
}

func (c *AdminController) ApplyHTTP(g *echo.Group) {
	g.GET("/admin/*/", c.Home, c.token)
	g.GET("/admin/callback/", c.Callback)
	g.POST("/admin/logout/", c.Logout, c.token)
}
