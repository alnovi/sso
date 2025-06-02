package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestMiddlewareTrace() {
	mTracer := middleware.Tracer()
	mErr := func(_ echo.HandlerFunc) echo.HandlerFunc {
		return func(_ echo.Context) error {
			return echo.ErrUnauthorized
		}
	}

	testCases := []struct {
		name       string
		mds        []echo.MiddlewareFunc
		traceId    string
		target     string
		expTraceId string
	}{
		{
			name:       "Trace 1",
			mds:        []echo.MiddlewareFunc{mTracer},
			traceId:    uuid.NewString(),
			target:     "/",
			expTraceId: "00000000000000000000000000000000",
		}, {
			name:       "Trace 2",
			mds:        []echo.MiddlewareFunc{mTracer},
			traceId:    uuid.NewString(),
			target:     "/favicon.png",
			expTraceId: "",
		}, {
			name:       "Trace 3",
			mds:        []echo.MiddlewareFunc{mTracer},
			traceId:    uuid.NewString(),
			target:     "/assets/file",
			expTraceId: "",
		}, {
			name:       "Trace 4",
			mds:        []echo.MiddlewareFunc{mTracer},
			traceId:    uuid.NewString(),
			target:     "/public/file",
			expTraceId: "",
		}, {
			name:       "Trace 5",
			mds:        []echo.MiddlewareFunc{mTracer, mErr},
			traceId:    uuid.NewString(),
			target:     "/*/",
			expTraceId: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Trace-ID", tc.traceId)
			rec := httptest.NewRecorder()
			ctx := s.app.HttpServer.NewContext(req, rec)
			ctx.SetPath(tc.target)

			if err := s.sendToMiddleware(ctx, tc.mds...); err != nil {
				s.Assert().ErrorContains(err, echo.ErrUnauthorized.Error(), MsgNotAssertError)
			}

			s.Assert().Equal(tc.expTraceId, rec.Header().Get("X-Trace-ID"), "not assert trace id")
		})
	}
}
