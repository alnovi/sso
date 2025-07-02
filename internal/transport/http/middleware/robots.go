package middleware

import "github.com/labstack/echo/v4"

func Robots() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			e.Response().Header().Set("X-Robots-Tag", "noindex, nofollow")
			return next(e)
		}
	}
}
