package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthClient() {
	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success get client",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			expCode: http.StatusOK,
			expBody: "",
			expErr:  "",
		},
		{
			name: "Response type invalid",
			query: map[string]string{
				"response_type": "invalid",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "response type invalid",
		},
		{
			name: "Client id invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     "invalid",
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "client not found",
		},
		{
			name: "Redirect uri invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  "invalid",
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "redirect url invalid",
		},
	}

	controller := oauth.NewAuthorizeController(s.app.Provider.OAuth(), s.app.Provider.Cookie())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)

			req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(controller.Client, c); err != nil {
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
