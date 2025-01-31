package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthTokenByCode() {
	session, err := s.app.Provider.Session().Create(
		context.Background(),
		s.app.Provider.Config().Admin.Id,
		TestIP,
		TestAgent,
	)
	s.Require().NoError(err)

	code, err := s.app.Provider.Token().CodeToken(
		context.Background(),
		session.Id,
		s.app.Provider.Config().Client.Id,
		s.app.Provider.Config().Admin.Id,
	)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expErr  string
	}{
		{
			name: "Successful token by code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Token by code is used",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "code incorrect",
		},
		{
			name: "Query empty code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          "",
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "code incorrect",
		},
		{
			name: "Query empty client_id",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     "",
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "client not found",
		},
		{
			name: "Query empty client_secret",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": "",
			},
			expCode: http.StatusBadRequest,
			expErr:  "client not found",
		},
		{
			name: "Grant type invalid",
			query: map[string]string{
				"grant_type":    "invalid",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "invalid grant_type",
		},
		{
			name: "Client id invalid",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     "invalid",
				"client_secret": s.app.Provider.Config().Client.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "client not found",
		},
		{
			name: "Client secret invalid",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().Client.Id,
				"client_secret": "invalid",
			},
			expCode: http.StatusBadRequest,
			expErr:  "client not found",
		},
		{
			name: "Code for another client",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"code":          code.Hash,
				"client_id":     s.app.Provider.Config().TestClient.Id,
				"client_secret": s.app.Provider.Config().TestClient.Secret,
			},
			expCode: http.StatusBadRequest,
			expErr:  "code incorrect",
		},
	}

	controller := oauth.NewTokenController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/?"+s.buildQuery(tc.query), nil)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(controller.Token, c); err != nil {
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
