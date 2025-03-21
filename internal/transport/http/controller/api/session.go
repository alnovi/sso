package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/storage"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type SessionController struct {
	controller.BaseController
	sessions *storage.Sessions
}

func NewSessionController(sessions *storage.Sessions) *SessionController {
	return &SessionController{sessions: sessions}
}

func (c *SessionController) List(e echo.Context) error {
	userSessionId := c.MustSessionId(e)

	sessions, err := c.sessions.List(e.Request().Context())
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewSessionsUser(sessions, userSessionId))
}

func (c *SessionController) Get(e echo.Context) error {
	userSessionId := c.MustSessionId(e)

	session, err := c.sessions.GetById(e.Request().Context(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewSessionUser(session, userSessionId))
}

func (c *SessionController) Delete(e echo.Context) error {
	if c.MustSessionId(e) == e.Param("id") {
		return echo.NewHTTPError(http.StatusBadRequest, "вы не можете удалить текущую сессию")
	}

	if err := c.sessions.DeleteById(e.Request().Context(), e.Param("id")); err != nil {
		return err
	}

	return e.NoContent(http.StatusOK)
}

func (c *SessionController) ApplyHTTP(g *echo.Group) error {
	g.GET("/sessions/", c.List)
	g.GET("/sessions/:id/", c.Get)
	g.DELETE("/sessions/:id/", c.Delete)
	return nil
}
