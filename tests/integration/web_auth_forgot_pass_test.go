package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (s *TestSuite) TestWebAuthForgotPasswordPage() {
	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expErr  error
	}{
		{
			name: "Success",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			expCode: http.StatusOK,
			expErr:  nil,
		},
		{
			name:    "Empty query",
			query:   map[string]string{},
			expCode: http.StatusFound,
			expErr:  nil,
		},
		{
			name: "Invalid response type",
			query: map[string]string{
				"response_type": "invalid",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrUnsupportedGrantType,
		},
		{
			name: "Invalid client",
			query: map[string]string{
				"response_type": "code",
				"client_id":     "invalid",
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrClientNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/?"+s.BuildQuery(tc.query), nil)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err := s.SendToServer(s.App.Provider.WebAuth().ForgotPasswordPage, c); err != nil {
				s.Assert().ErrorAs(err, &tc.expErr, "not assert error") //nolint:gosec
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}

func (s *TestSuite) TestWebAuthForgotPassword() {
	testCases := []struct {
		name    string
		mime    string
		query   map[string]string
		form    map[string]string
		expCode int
		expErr  error
	}{
		{
			name: "Success form",
			mime: echo.MIMEApplicationForm,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			form: map[string]string{
				"login": s.App.Config().User.AdminEmail,
			},
			expCode: http.StatusOK,
			expErr:  nil,
		},
		{
			name: "Success json",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			form: map[string]string{
				"login": s.App.Config().User.AdminEmail,
			},
			expCode: http.StatusOK,
			expErr:  nil,
		},
		{
			name: "Invalid response type",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "invalid",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			form: map[string]string{
				"login": s.App.Config().User.AdminEmail,
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrUnsupportedGrantType,
		},
		{
			name: "Invalid client",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     "invalid",
			},
			form: map[string]string{
				"login": s.App.Config().User.AdminEmail,
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrClientNotFound,
		},
		{
			name: "Invalid email format",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			form: map[string]string{
				"login": "invalid",
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "Invalid login",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Config().Client.ProfileID,
			},
			form: map[string]string{
				"login": "invalid@example.com",
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  &validator.ValidateError{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.BuildQuery(tc.query)
			data := s.BuildFormData(tc.mime, tc.form)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, strings.NewReader(data))
			req.Header.Set(echo.HeaderContentType, tc.mime)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err := s.SendToServer(s.App.WebAuth().ForgotPassword, c); err != nil {
				s.Assert().NotNil(tc.expErr, fmt.Sprintf("err is not nil: %s", err))
				s.Assert().ErrorAs(err, &tc.expErr, fmt.Sprintf("not assert error: %s", err)) //nolint:gosec
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}
