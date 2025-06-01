package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/alnovi/sso/internal/helper"
)

const (
	HeaderXTraceID = "X-Trace-ID"
)

func Tracer() echo.MiddlewareFunc {
	getSpanName := func(e echo.Context) string {
		spanName := fmt.Sprintf("%s - not found", e.Request().Method)
		if path := e.Path(); path != "" {
			spanName = fmt.Sprintf("%s - %s", e.Request().Method, path)
		}
		return spanName
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			var err error

			isFavicon := strings.Contains(e.Path(), "favicon")
			isAssets := strings.Contains(e.Path(), "assets")
			isPublic := strings.Contains(e.Path(), "public")
			isHttp := strings.Contains(e.Path(), "/*/")

			if isFavicon || isAssets || isPublic || isHttp {
				return next(e)
			}

			attr := []attribute.KeyValue{
				attribute.String("request.id", e.Response().Header().Get(echo.HeaderXRequestID)),
				semconv.HTTPRoute(e.Path()),
			}

			for _, param := range e.ParamNames() {
				attr = append(attr, attribute.String(fmt.Sprintf("http.route.param.%s", param), e.Param(param)))
			}

			ctx := otel.GetTextMapPropagator().Extract(e.Request().Context(), propagation.HeaderCarrier(e.Request().Header))
			ctx, span := helper.SpanStart(ctx, getSpanName(e), trace.WithAttributes(attr...))
			defer span.End()

			e.Response().Header().Set(HeaderXTraceID, span.SpanContext().TraceID().String())
			e.SetRequest(e.Request().WithContext(ctx))

			if err = next(e); err != nil {
				helper.SpanError(span, err)
			}

			return err
		}
	}
}
