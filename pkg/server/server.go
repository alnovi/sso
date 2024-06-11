package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Option func(e *echo.Echo)

type Controller interface {
	Route(e *echo.Group)
}

type Server struct {
	*echo.Echo
}

func New(options ...Option) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.AddTrailingSlash())

	for _, option := range options {
		option(e)
	}

	return &Server{Echo: e}
}

func (s *Server) ApplyController(prefix string, cs []Controller, ms ...echo.MiddlewareFunc) {
	g := s.Group(prefix, ms...)
	for _, c := range cs {
		c.Route(g)
	}
}

func WithErrorHandle(handle echo.HTTPErrorHandler) Option {
	return func(e *echo.Echo) {
		e.HTTPErrorHandler = handle
	}
}

func WithValidate(validator echo.Validator) Option {
	return func(e *echo.Echo) {
		e.Validator = validator
	}
}

func WithRender(render echo.Renderer) Option {
	return func(e *echo.Echo) {
		e.Renderer = render
	}
}
