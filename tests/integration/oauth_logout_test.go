package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthLogout() {
	session, err := s.app.Provider.Session().Create(context.Background(), s.app.Provider.Config().Admin.Id, "", "")
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		cookie  *http.Cookie
		expCode int
		expErr  string
	}{
		{
			name:    "Success session id",
			cookie:  s.app.Provider.Cookie().SessionId(session.Id),
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name:    "Invalid session id",
			cookie:  s.app.Provider.Cookie().SessionId("invalid"),
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name:    "Empty session id",
			cookie:  nil,
			expCode: http.StatusOK,
			expErr:  "",
		},
	}

	controller := oauth.NewProfileController(s.app.Provider.OAuth(), s.app.Provider.Cookie())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(controller.Logout, c); err != nil {
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
