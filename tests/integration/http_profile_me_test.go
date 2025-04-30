package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpProfileMe() {
	session, _, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
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
			name: "Success",
			headers: map[string]string{
				"User-Agent": TestAgent,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusOK,
			expBody: s.config().UAdmin.Id,
		},
		{
			name:    "Unauthorized",
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized: session not found",
		},
		{
			name: "Invalid user agent",
			headers: map[string]string{
				"User-Agent": "invalid",
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized: session not found: agent not attempted",
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Me, c, mdw); err != nil {
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
