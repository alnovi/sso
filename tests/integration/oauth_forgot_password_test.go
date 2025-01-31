package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthForgotPassword() {
	testCases := []struct {
		name    string
		query   map[string]string
		data    map[string]any
		headers map[string]string
		expCode int
		expErr  string
	}{
		{
			name: "Successful forgot password json",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": s.app.Provider.Config().Admin.Email,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Successful forgot password form",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": s.app.Provider.Config().Admin.Email,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationForm,
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Response type invalid",
			query: map[string]string{
				"response_type": "invalid",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": s.app.Provider.Config().Admin.Email,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "response type invalid",
		},
		{
			name: "User login is empty",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": "",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User login invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": "invalid",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User login not found",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login": "invalid@example.com",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  "Unprocessable Entity",
		},
	}

	controller := oauth.NewPasswordController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, strings.NewReader(data))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(controller.ForgotPassword, c); err != nil {
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
