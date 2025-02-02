package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BaseController struct{}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (c *BaseController) BindValidate(e echo.Context, trg any) error {
	if err := e.Bind(trg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request").SetInternal(err)
	}

	if err := e.Validate(trg); err != nil {
		return err
	}

	return nil
}

func (c *BaseController) MustUserId(e echo.Context) string {
	userId, ok := e.Get("user_id").(string)
	if !ok {
		panic(errors.New("context user id not found"))
	}
	return userId
}

func (c *BaseController) MustClientId(e echo.Context) string {
	clientId, ok := e.Get("client_id").(string)
	if !ok {
		panic(errors.New("context client id not found"))
	}
	return clientId
}

func (c *BaseController) MustUserRole(e echo.Context) string {
	role, ok := e.Get("user_role").(string)
	if !ok {
		panic(errors.New("context user role not found"))
	}
	return role
}
