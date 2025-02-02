package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transaport/http/middleware"
)

func (s *TestSuite) TestOauthProfile() {
	_, token, err := s.app.Provider.JWT().GenerateAccessToken(
		s.app.Provider.Config().Client.Id,
		s.app.Provider.Config().Admin.Id,
		TestRoleAdmin,
	)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		expCode int
		expErr  string
	}{
		{
			name: "Success with token",
			headers: map[string]string{
				echo.HeaderAuthorization: token,
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Success with bearer token",
			headers: map[string]string{
				echo.HeaderAuthorization: fmt.Sprintf("Bearer %s", token),
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Empty token",
			headers: map[string]string{
				echo.HeaderAuthorization: "",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "unauthenticated",
		},
		{
			name: "Invalid token",
			headers: map[string]string{
				echo.HeaderAuthorization: "invalid",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "unauthenticated",
		},
	}

	middlewares := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.JWT()),
	}

	ctr := oauth.NewProfileController(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.User())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctr.Logout, c, middlewares...); err != nil {
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
