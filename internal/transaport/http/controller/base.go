package controller

import (
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
