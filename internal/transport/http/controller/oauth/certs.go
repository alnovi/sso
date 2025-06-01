package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/certs"
	"github.com/alnovi/sso/internal/transport/http/controller"
)

type CertsController struct {
	controller.BaseController
	certs *certs.Certs
}

func NewCertsController(certs *certs.Certs) *CertsController {
	return &CertsController{certs: certs}
}

func (c *CertsController) Certs(e echo.Context) error {
	jwk, err := c.certs.PublicJWK()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error generate JWK").SetInternal(err)
	}
	return e.JSON(http.StatusOK, jwk)
}

func (c *CertsController) ApplyHTTP(g *echo.Group) {
	g.GET("/certs/", c.Certs)
}
