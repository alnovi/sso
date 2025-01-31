package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/oauth"
	"github.com/alnovi/sso/internal/transaport/http/controller"
)

type PasswordController struct {
	*controller.BaseController
	oauth *oauth.OAuth
}

func NewPasswordController(oauth *oauth.OAuth) *PasswordController {
	return &PasswordController{oauth: oauth}
}

func (c *PasswordController) ForgotPassword(e echo.Context) error {
	ctx := e.Request().Context()

	if _, err := c.oauth.ResponseType(e.QueryParam("response_type")); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_ = ctx

	return nil
}

func (c *PasswordController) ResetPassword(e echo.Context) error {
	return nil
}

func (c *PasswordController) ApplyHTTP(g *echo.Group) error {
	g.POST("/forgot-password/", c.ForgotPassword)
	g.POST("/reset-password/", c.ResetPassword)
	return nil
}
