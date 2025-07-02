package api

import (
	"context"
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

// List         godoc
// @Id          SessionsList
// @Summary     Список сессий
// @Description Список сессий пользователей
// @Tags        Api Sessions
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Success 200 {object} []response.SessionUser "Информация о сессиях пользователей"
// @Failure 403
// @Router      /api/sessions [get]
func (c *SessionController) List(e echo.Context) error {
	userSessionId := c.MustSessionId(e)

	sessions, err := c.sessions.List(context.Background())
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewSessionsUser(sessions, userSessionId))
}

// Get          godoc
// @Id          SessionsGet
// @Summary     Сессия
// @Description Информация о сессии пользователя
// @Tags        Api Sessions
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор сессии"
// @Success 200 {object} response.SessionUser "Информация о сессии пользователя"
// @Failure 403
// @Failure 404
// @Router      /api/sessions/{id} [get]
func (c *SessionController) Get(e echo.Context) error {
	userSessionId := c.MustSessionId(e)

	session, err := c.sessions.GetById(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewSessionUser(session, userSessionId))
}

// Delete       godoc
// @Id          SessionsDelete
// @Summary     Удаление сессии
// @Description Удаление сессии пользователя
// @Tags        Api Sessions
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор сессии"
// @Success 200
// @Failure 400
// @Failure 403
// @Router      /api/sessions/{id} [delete]
func (c *SessionController) Delete(e echo.Context) error {
	if c.MustSessionId(e) == e.Param("id") {
		return echo.NewHTTPError(http.StatusBadRequest, "вы не можете удалить текущую сессию")
	}

	if err := c.sessions.DeleteById(context.Background(), e.Param("id")); err != nil {
		return err
	}

	return e.NoContent(http.StatusOK)
}

func (c *SessionController) ApplyHTTP(g *echo.Group) {
	g.GET("/sessions/", c.List)
	g.GET("/sessions/:id/", c.Get)
	g.DELETE("/sessions/:id/", c.Delete)
}
