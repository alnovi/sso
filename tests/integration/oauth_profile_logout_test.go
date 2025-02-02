package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transaport/http/middleware"
)

func (s *TestSuite) TestOauthProfileLogout() {
	_, token, err := s.app.Provider.JWT().GenerateAccessToken(
		s.app.Provider.Config().Client.Id,
		s.app.Provider.Config().Admin.Id,
		TestRoleAdmin,
	)
	s.Require().NoError(err)

	session, err := s.app.Provider.Session().Create(context.Background(), s.app.Provider.Config().Admin.Id, TestIP, TestAgent)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		cookies []*http.Cookie
		expCode int
		expErr  string
	}{
		{
			name: "Success session id",
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id),
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Invalid session id",
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId("invalid"),
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Unauthorized session id",
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id),
			},
			expCode: http.StatusUnauthorized,
			expErr:  "unauthenticated",
		},
		{
			name: "Empty session id",
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
	}

	middlewares := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.JWT()),
	}

	controller := oauth.NewProfileController(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.User())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			for _, cookie := range tc.cookies {
				req.AddCookie(cookie)
			}

			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(controller.Logout, c, middlewares...); err != nil {
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
