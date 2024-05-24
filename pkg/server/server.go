package server

import (
	"errors"
	"log/slog"
	"net/http"

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

func WithLogger(logger *slog.Logger) Option {
	return func(e *echo.Echo) {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			//HandleError:      true,
			LogRemoteIP:     true,
			LogUserAgent:    true,
			LogLatency:      true,
			LogMethod:       true,
			LogURI:          true,
			LogStatus:       true,
			LogError:        true,
			LogResponseSize: true,
			LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
				log := logger.With(
					slog.String("module", "http-request"),
					slog.String("method", values.Method),
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
					slog.Int64("response_size", values.ResponseSize),
					slog.Duration("latency", values.Latency),
					slog.String("ip", values.RemoteIP),
					slog.String("agent", values.UserAgent),
				)

				var echoHttpError *echo.HTTPError
				if errors.As(values.Error, &echoHttpError) {
					if echoHttpError.Internal != nil {
						log = log.With("error", echoHttpError.Internal.Error())
					}
					log.Error(echoHttpError.Message.(string))
				} else if values.Error != nil {
					log.Error("Internal Server Error", "error", values.Error.Error())
				} else {
					log.Info(http.StatusText(values.Status))
				}

				return nil
			},
		}))
	}
}

func WithCors(origin string) Option {
	return func(e *echo.Echo) {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			Skipper:          middleware.DefaultSkipper,
			AllowOrigins:     []string{origin},
			AllowMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowCredentials: true,
		}))
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
