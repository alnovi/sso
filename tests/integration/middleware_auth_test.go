package integration

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestMiddlewareAuth() {
	authAdmin := middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret)
	authTest := middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), TestClient.Id, TestClient.Secret)

	session1, access1, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	session2, _, refresh2, err := s.accessTokens(TestClient.Id, TestUser.Id, TestRole, token.WithAccessExpiresAt(time.Now()))
	s.Require().NoError(err)

	session3, _, refresh3, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin, token.WithAccessExpiresAt(time.Now()))
	s.Require().NoError(err)

	testCases := []struct {
		name       string
		handler    echo.MiddlewareFunc
		headers    map[string]string
		cookies    []*http.Cookie
		expCode    int
		expSession any
		expClient  any
		expUser    any
		expRole    any
		expErr     string
	}{
		{
			name:    "Success access header",
			handler: authAdmin,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access1.Hash,
			},
			expCode:    http.StatusOK,
			expSession: session1.Id,
			expClient:  s.config().CAdmin.Id,
			expUser:    s.config().UAdmin.Id,
			expRole:    entity.RoleAdmin,
		}, {
			name:    "Success access cookie",
			handler: authAdmin,
			headers: map[string]string{
				"User-Agent":   TestAgent,
				"Content-Type": "application/json",
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().AccessToken(access1),
			},
			expCode:    http.StatusOK,
			expSession: session1.Id,
			expClient:  s.config().CAdmin.Id,
			expUser:    s.config().UAdmin.Id,
			expRole:    entity.RoleAdmin,
		}, {
			name:    "Success refresh header",
			handler: authTest,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": "invalid-token",
				"Refresh-Token": refresh2.Hash,
			},
			expCode:    http.StatusOK,
			expSession: session2.Id,
			expClient:  TestClient.Id,
			expUser:    TestUser.Id,
			expRole:    TestRole,
		}, {
			name:    "Success refresh cookie",
			handler: authAdmin,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": "invalid-token",
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().RefreshToken(refresh3),
			},
			expCode:    http.StatusOK,
			expSession: session3.Id,
			expClient:  s.config().CAdmin.Id,
			expUser:    s.config().UAdmin.Id,
			expRole:    entity.RoleAdmin,
		}, {
			name:    "Unauthorized",
			handler: authAdmin,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": "invalid-token",
				"Refresh-Token": "invalid-token",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		}, {
			name:    "Unauthorized is user refresh",
			handler: authTest,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": "invalid-token",
				"Refresh-Token": refresh2.Hash,
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()
			ctx := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToMiddleware(ctx, tc.handler); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)

			s.Assert().Equal(tc.expSession, ctx.Get(controller.CtxSessionId), "not assert session")
			s.Assert().Equal(tc.expClient, ctx.Get(controller.CtxClientId), "not assert client id")
			s.Assert().Equal(tc.expUser, ctx.Get(controller.CtxUserId), "not assert user id")
			s.Assert().Equal(tc.expRole, ctx.Get(controller.CtxUserRole), "not assert user role")
		})
	}
}
