package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpOAuthForgotPassword() {
	testCases := []struct {
		name    string
		query   map[string]string
		headers map[string]string
		data    map[string]any
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success",
			query: map[string]string{
				"client_id":    s.config().CAdmin.Id,
				"redirect_uri": s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login": s.config().UAdmin.Email,
			},
			expCode: http.StatusOK,
			expBody: "Ссылка для смены пароля отправлена на электронную почту",
		}, {
			name: "Empty login",
			query: map[string]string{
				"client_id":    s.config().CAdmin.Id,
				"redirect_uri": s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "login обязательное поле",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Invalid login",
			query: map[string]string{
				"client_id":    s.config().CAdmin.Id,
				"redirect_uri": s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    "invalid",
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "login должен быть email адресом",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Not found user",
			query: map[string]string{
				"client_id":    s.config().CAdmin.Id,
				"redirect_uri": s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login": "invalid@example.com",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "пользователь не найден",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Invalid client",
			query: map[string]string{
				"client_id":    "invalid",
				"redirect_uri": s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login": s.config().UAdmin.Email,
			},
			expCode: http.StatusBadRequest,
			expBody: "Клиент не найден",
			expErr:  "client not found",
		}, {
			name: "Invalid redirect URI",
			query: map[string]string{
				"client_id":    s.config().CAdmin.Id,
				"redirect_uri": "invalid",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login": s.config().UAdmin.Email,
			},
			expCode: http.StatusBadRequest,
			expBody: "Не валидный redirect-uri",
			expErr:  "invalid redirect uri",
		},
	}

	ctrl := oauth.NewPasswordController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(ctrl.ForgotPassword, c, middleware.TrailingSlash()); err != nil {
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
