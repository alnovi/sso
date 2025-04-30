package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	CtxSessionId = "session_id"
	CtxClientId  = "client_id"
	CtxUserId    = "user_id"
	CtxUserRole  = "user_role"
)

type BaseController struct{}

func (c *BaseController) SessionId(e echo.Context) (string, bool) {
	val, ok := e.Get(CtxSessionId).(string)
	return val, ok
}

func (c *BaseController) MustSessionId(e echo.Context) string {
	val, ok := c.SessionId(e)
	if !ok {
		panic("session id not found")
	}
	return val
}

func (c *BaseController) UserId(e echo.Context) (string, bool) {
	val, ok := e.Get(CtxUserId).(string)
	return val, ok
}

func (c *BaseController) MustUserId(e echo.Context) string {
	val, ok := c.UserId(e)
	if !ok {
		panic("user id not found")
	}
	return val
}

func (c *BaseController) BindValidate(e echo.Context, dst any) error {
	if err := e.Bind(dst); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	return e.Validate(dst)
}
