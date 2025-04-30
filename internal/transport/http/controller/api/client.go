package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/storage"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
	"github.com/alnovi/sso/pkg/validator"
)

type ClientController struct {
	controller.BaseController
	clients *storage.Clients
}

func NewClientController(clients *storage.Clients) *ClientController {
	return &ClientController{clients: clients}
}

func (c *ClientController) List(e echo.Context) error {
	clients, err := c.clients.All(e.Request().Context())
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClients(clients))
}

func (c *ClientController) Get(e echo.Context) error {
	client, err := c.clients.GetById(e.Request().Context(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClient(client))
}

func (c *ClientController) Create(e echo.Context) error {
	req := new(request.CreateClient)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := storage.InputClientCreate{
		Id:       req.Id,
		Name:     req.Name,
		Icon:     req.Icon,
		Callback: req.Callback,
		Secret:   req.Secret,
	}

	client, err := c.clients.Create(e.Request().Context(), inp)
	if err != nil {
		if errors.Is(err, storage.ErrClientIdExists) {
			return validator.NewValidateErrorWithMessage("id", "Такое значение уже занято")
		}
		return err
	}

	return e.JSON(http.StatusOK, response.NewClient(client))
}

func (c *ClientController) Update(e echo.Context) error {
	req := new(request.UpdateClient)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := storage.InputClientUpdate{
		Id:       e.Param("id"),
		Name:     req.Name,
		Icon:     req.Icon,
		Callback: req.Callback,
		Secret:   req.Secret,
	}

	client, err := c.clients.Update(e.Request().Context(), inp)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, response.NewClient(client))
}

func (c *ClientController) Delete(e echo.Context) error {
	client, err := c.clients.Delete(e.Request().Context(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClient(client))
}

func (c *ClientController) Restore(e echo.Context) error {
	client, err := c.clients.Restore(e.Request().Context(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClient(client))
}

func (c *ClientController) ApplyHTTP(g *echo.Group) error {
	g.GET("/clients/", c.List)
	g.GET("/clients/:id/", c.Get)
	g.POST("/clients/", c.Create)
	g.PUT("/clients/:id/", c.Update)
	g.DELETE("/clients/:id/", c.Delete)
	g.POST("/clients/:id/restore/", c.Restore)
	return nil
}
