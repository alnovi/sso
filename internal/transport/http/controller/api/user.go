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

// List         godoc
// @Id          UsersList
// @Summary     Список пользователей
// @Description Список пользователей
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Success 200 {object} []response.User "Информация о пользователях"
// @Failure 403
// @Router      /api/users [get]
func (c *UserController) List(e echo.Context) error {
	users, err := c.users.All(context.Background())
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUsers(users))
}

// Get          godoc
// @Id          UsersGet
// @Summary     Пользователь
// @Description Информация о пользователе
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор пользователя"
// @Success 200 {object} response.User "Информация о пользователе"
// @Failure 403
// @Failure 404
// @Router      /api/users/{id} [get]
func (c *UserController) Get(e echo.Context) error {
	user, err := c.users.GetById(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

// Clients      godoc
// @Id          UsersClients
// @Summary     Список клиентов пользователя
// @Description Список приложений и права доступа для пользователя
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Success 200 {object} []response.ClientRole "Информация о клиентах и правах доступа"
// @Failure 403
// @Router      /api/users/{id}/clients [get]
func (c *UserController) Clients(e echo.Context) error {
	clientRole, err := c.roles.ClientRoleByUserId(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewClientsRoles(clientRole))
}

// Create       godoc
// @Id          UsersCreate
// @Summary     Добавление пользователя
// @Description Добавление нового пользователя
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       request body request.CreateUser true "Данные пользователя"
// @Success 200 {object} response.User "Информация о пользователе"
// @Failure 403
// @Failure 422
// @Router      /api/users [post]
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

// Update       godoc
// @Id          UsersUpdate
// @Summary     Изменение пользователя
// @Description Изменение пользователя
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор пользователя"
// @Param       request body request.UpdateUser true "Данные пользователя"
// @Success 200 {object} response.User "Информация о пользователе"
// @Failure 403
// @Failure 404
// @Failure 422
// @Router      /api/users/{id} [put]
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

// Delete       godoc
// @Id          UsersDelete
// @Summary     Удаление пользователя
// @Description Удаление пользователя
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор пользователя"
// @Success 200 {object} response.User "Информация о пользователе"
// @Failure 403
// @Failure 404
// @Router      /api/users/{id} [delete]
func (c *UserController) Delete(e echo.Context) error {
	user, err := c.users.Delete(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

// Restore      godoc
// @Id          UsersRestore
// @Summary     Восстановление пользователя
// @Description Восстановление удаленного пользователя
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       id query string true "Идентификатор пользователя"
// @Success 200 {object} response.User "Информация о пользователе"
// @Failure 403
// @Failure 404
// @Router      /api/users/{id}/restore [post]
func (c *UserController) Restore(e echo.Context) error {
	user, err := c.users.Restore(context.Background(), e.Param("id"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, response.NewUser(user))
}

// UpdateRole   godoc
// @Id          UsersUpdateRole
// @Summary     Изменение роли пользователя
// @Description Изменение роли доступа пользователя в приложении
// @Tags        Api Users
// @Accept      json
// @Produce     json
// @Security    JWT-Access
// @Param       uid query string true "Идентификатор пользователя"
// @Param       сid query string true "Идентификатор клиента"
// @Param       request body request.UpdateUserRole true "Роль пользователя"
// @Success 200
// @Failure 403
// @Failure 404
// @Failure 422
// @Router      /api/users/{uid}/clients/{cid} [post]
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
