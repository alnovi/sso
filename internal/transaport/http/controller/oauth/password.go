package oauth

import (
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

func (c *PasswordController) ForgotPassword(ctx echo.Context) error {
	return nil
}

func (c *PasswordController) ResetPassword(ctx echo.Context) error {
	return nil
}

func (c *PasswordController) ApplyHTTP(g *echo.Group) error {
	g.POST("/forgot-password/", c.ForgotPassword)
	g.POST("/reset-password/", c.ResetPassword)
	return nil
}
