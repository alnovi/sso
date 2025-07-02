package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/service/certs"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/response"
)

type CertsController struct {
	controller.BaseController
	certs *certs.Certs
}

func NewCertsController(certs *certs.Certs) *CertsController {
	return &CertsController{certs: certs}
}

// Certs        godoc
// @Id          OAuthCerts
// @Summary     Публичный ключ
// @Description Публичный ключ для проверки токена
// @Tags        OAuth
// @Accept      json
// @Produce     json
// @Success 200 {object} response.JWK "Json web key"
// @Failure 500
// @Router      /oauth/certs [get]
func (c *CertsController) Certs(e echo.Context) error {
	jwk, err := c.certs.PublicJWK()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error generate JWK").SetInternal(err)
	}
	return e.JSON(http.StatusOK, response.NewJWK(jwk))
}

func (c *CertsController) ApplyHTTP(g *echo.Group) {
	g.GET("/certs/", c.Certs)
}
