package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RequestLogger(logger *slog.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	config := middleware.RequestLoggerConfig{
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
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
					log.With("internal", echoHttpError.Internal.Error())
				}
				log.Error(echoHttpError.Message.(string))
			} else if values.Error != nil {
				log.Error("Internal Server Error", "internal", values.Error.Error())
			} else {
				log.Info(http.StatusText(values.Status))
			}

			return nil
		},
	}

	return middleware.RequestLoggerWithConfig(config)
}
