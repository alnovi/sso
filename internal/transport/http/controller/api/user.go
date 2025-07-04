package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/alnovi/gomon/validator"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/storage"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/request"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type UserController struct {
	controller.BaseController
	users *storage.Users
	roles *storage.Roles
}

func NewUserController(users *storage.Users, roles *storage.Roles) *UserController {
	return &UserController{users: users, roles: roles}
}

func (c *UserController) List(e echo.Context) error {
	users, err := c.users.All(context.Background())
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUsers(users))
}

func (c *UserController) Get(e echo.Context) error {
	user, err := c.users.GetById(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

func (c *UserController) Clients(e echo.Context) error {
	clientRole, err := c.roles.ClientRoleByUserId(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClientsRoles(clientRole))
}

func (c *UserController) Create(e echo.Context) error {
	req := new(request.CreateUser)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := storage.InputUserCreate{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := c.users.Create(context.Background(), inp)
	if err != nil {
		if errors.Is(err, storage.ErrUserEmailExists) {
			return validator.NewValidateErrorWithMessage("email", "Такое значение уже занято")
		}
		return err
	}

	return e.JSON(http.StatusOK, response.NewUser(user))
}

func (c *UserController) Update(e echo.Context) error {
	req := new(request.UpdateUser)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	inp := storage.InputUserUpdate{
		Id:       e.Param("id"),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := c.users.Update(context.Background(), inp)
	if err != nil {
		if errors.Is(err, storage.ErrUserEmailExists) {
			return validator.NewValidateErrorWithMessage("email", "Такое значение уже занято")
		}
		return err
	}

	return e.JSON(http.StatusOK, response.NewUser(user))
}

func (c *UserController) Delete(e echo.Context) error {
	user, err := c.users.Delete(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

func (c *UserController) Restore(e echo.Context) error {
	user, err := c.users.Restore(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

func (c *UserController) UpdateRole(e echo.Context) error {
	ctx := context.Background()
	clientId := e.Param("cid")
	userId := e.Param("uid")

	req := new(request.UpdateUserRole)

	if err := c.BindValidate(e, req); err != nil {
		return err
	}

	if err := c.roles.Update(ctx, clientId, userId, req.Role); err != nil {
		return err
	}

	return e.NoContent(http.StatusOK)
}

func (c *UserController) ApplyHTTP(g *echo.Group) {
	g.GET("/users/", c.List)
	g.GET("/users/:id/", c.Get)
	g.GET("/users/:id/clients/", c.Clients)
	g.POST("/users/", c.Create)
	g.PUT("/users/:id/", c.Update)
	g.DELETE("/users/:id/", c.Delete)
	g.POST("/users/:id/restore/", c.Restore)
	g.POST("/users/:uid/clients/:cid/", c.UpdateRole)
}
