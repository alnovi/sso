package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpOauthFormResetPassword() {
	query := s.buildQuery(map[string]string{
		"client_id":    s.config().CAdmin.Id,
		"redirect_uri": s.config().CAdmin.Callback,
	})

	token, err := s.app.Provider.Token().ForgotPasswordToken(
		context.Background(),
		s.config().CAdmin.Id,
		s.config().UAdmin.Id,
		query,
		TestIP,
		TestAgent,
	)

	s.Require().NoError(err)

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success",
			query: map[string]string{
				"hash": token.Hash,
			},
			expCode: http.StatusOK,
			expBody: `<div id="app"></div>`,
		}, {
			name:    "Hash is empty",
			query:   map[string]string{},
			expCode: http.StatusBadRequest,
			expBody: `Токен не найден`,
			expErr:  "token not found",
		}, {
			name: "Hash is invalid",
			query: map[string]string{
				"hash": "invalid",
			},
			expCode: http.StatusBadRequest,
			expBody: `Токен не найден`,
			expErr:  "token not found",
		},
	}

	ctrl := oauth.NewPasswordController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/?"+s.buildQuery(tc.query), nil)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.FormReset, c, middleware.TrailingSlash()); err != nil {
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
