package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
)

func RoleWeight(weight int) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			userRole, ok := e.Get(controller.CtxUserRole).(string)
			if !ok {
				return echo.ErrForbidden
			}

			if entity.RoleMap[userRole] < weight {
				return echo.ErrForbidden
			}

			return next(e)
		}
	}
}
