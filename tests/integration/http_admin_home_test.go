package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpAdminHome() {
	session, access, refresh, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		cookies []*http.Cookie
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success by access",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusOK,
			expBody: "SSO | Admin",
		}, {
			name: "Success by refresh",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Authorization": fmt.Sprintf("Bearer %s", "invalid"),
				"Refresh-Token": refresh.Hash,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusOK,
			expBody: "SSO | Admin",
		}, {
			name: "Unauthorized",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Authorization": fmt.Sprintf("Bearer %s", "invalid"),
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusFound,
			expBody: "",
		},
	}

	mdw := middleware.Token(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret)
	ctrl := controller.NewAdminController(s.app.Provider.Admin(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Home, c, mdw); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Contains(rec.Body.String(), tc.expBody, MsgNotAssertBody)

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
