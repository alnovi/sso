package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthResetPassword() {
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", s.app.Provider.Config().Client.Id)
	query.Set("redirect_uri", s.app.Provider.Config().Client.Host)

	reset, err := s.app.Provider.Token().ForgotPasswordToken(
		context.Background(),
		s.app.Provider.Config().Client.Id,
		s.app.Provider.Config().Admin.Id,
		query.Encode(),
	)
	s.Require().NoError(err)

	reset2, err := s.app.Provider.Token().ForgotPasswordToken(
		context.Background(),
		s.app.Provider.Config().Client.Id,
		s.app.Provider.Config().Admin.Id,
		query.Encode(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		query   map[string]string
		data    map[string]any
		headers map[string]string
		expCode int
		expErr  string
	}{
		{
			name: "Successful reset password json",
			query: map[string]string{
				"hash": reset.Hash,
			},
			data: map[string]any{
				"password": s.app.Provider.Config().Admin.Password,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusFound,
			expErr:  "",
		},
		{
			name:  "Successful reset password ajax",
			query: map[string]string{"hash": reset2.Hash},
			data:  map[string]any{"password": s.app.Provider.Config().Admin.Password},
			headers: map[string]string{
				"Content-Type":     echo.MIMEApplicationJSON,
				"X-Requested-With": "XMLHttpRequest",
			},
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name: "Reset password by hash is used",
			query: map[string]string{
				"hash": reset.Hash,
			},
			data: map[string]any{
				"password": s.app.Provider.Config().Admin.Password,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "token not found",
		},
		{
			name: "Hash empty",
			query: map[string]string{
				"hash": "",
			},
			data: map[string]any{
				"password": s.app.Provider.Config().Admin.Password,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "token not found",
		},
		{
			name: "Hash invalid",
			query: map[string]string{
				"hash": "invalid",
			},
			data: map[string]any{
				"password": s.app.Provider.Config().Admin.Password,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "token not found",
		},
		{
			name: "User password empty",
			query: map[string]string{
				"hash": reset.Hash,
			},
			data: map[string]any{
				"password": "",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User password invalid",
			query: map[string]string{
				"hash": reset.Hash,
			},
			data: map[string]any{
				"password": "123",
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
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/?"+s.buildQuery(tc.query), strings.NewReader(data))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(controller.ResetPassword, c); err != nil {
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
