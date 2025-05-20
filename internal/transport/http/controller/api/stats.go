package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/stats"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type StatsController struct {
	controller.BaseController
	stats *stats.Stats
}

func NewStatsController(stats *stats.Stats) *StatsController {
	return &StatsController{stats: stats}
}

func (c *StatsController) Stats(e echo.Context) error {
	var users, clients, sessions int
	var err error

	ctx := e.Request().Context()

	if users, err = c.stats.UserCount(ctx); err != nil {
		return err
	}

	if clients, err = c.stats.ClientCount(ctx); err != nil {
		return err
	}

	if sessions, err = c.stats.SessionCount(ctx); err != nil {
		return err
	}

	resp := &response.Stats{
		Users:    users,
		Clients:  clients,
		Sessions: sessions,
	}

	return e.JSON(http.StatusOK, resp)
}

func (c *StatsController) ApplyHTTP(g *echo.Group) {
	g.GET("/stats/", c.Stats)
}
