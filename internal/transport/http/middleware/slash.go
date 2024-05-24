package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TrailingSlash() func(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.AddTrailingSlash()
}
