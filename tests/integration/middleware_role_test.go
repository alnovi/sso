package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestMiddlewareRole() {
	roleAdmin := middleware.RoleWeight(entity.RoleAdminWeight)
	roleManager := middleware.RoleWeight(entity.RoleManagerWeight)
	roleUser := middleware.RoleWeight(entity.RoleUserWeight)

	testCases := []struct {
		name    string
		handler echo.MiddlewareFunc
		role    string
		expCode int
		expErr  string
	}{
		{
			name:    "Success admin admin",
			handler: roleAdmin,
			role:    entity.RoleAdmin,
			expCode: http.StatusOK,
		}, {
			name:    "Success admin manager",
			handler: roleManager,
			role:    entity.RoleAdmin,
			expCode: http.StatusOK,
		}, {
			name:    "Success admin user",
			handler: roleUser,
			role:    entity.RoleAdmin,
			expCode: http.StatusOK,
		}, {
			name:    "Success manager manager",
			handler: roleManager,
			role:    entity.RoleManager,
			expCode: http.StatusOK,
		}, {
			name:    "Success manager user",
			handler: roleUser,
			role:    entity.RoleManager,
			expCode: http.StatusOK,
		}, {
			name:    "Success user user",
			handler: roleUser,
			role:    entity.RoleUser,
			expCode: http.StatusOK,
		}, {
			name:    "Forbidden",
			handler: roleAdmin,
			role:    entity.RoleUser,
			expCode: http.StatusForbidden,
			expErr:  "Forbidden",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			ctx := s.app.HttpServer.NewContext(req, rec)
			ctx.Set(controller.CtxUserRole, tc.role)

			if err := s.sendToMiddleware(ctx, tc.handler); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
